package helper

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ratedemon/go-rest/api/middleware"
)

type HandleFunc func(ctx context.Context, req *http.Request) (interface{}, error)

type Route struct {
	Path, Method string
	HandleFunc   HandleFunc
}

type ApiHandler interface {
	RegisterRoutes() []Route
}

func HandleWrapper(f HandleFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		var res interface{}

		if userID, ok := r.Context().Value(middleware.UserIDKey).(int64); ok {
			res, err = f(contextWithUserID(context.Background(), userID), r)
		} else {
			res, err = f(r.Context(), r)
		}
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(err.Error())
		}

		body, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})
}
