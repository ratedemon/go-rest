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
		// fmt.Println(r.Context().Value(middleware.UserIDKey).(int64))
		if userID, ok := r.Context().Value(middleware.UserIDKey).(int64); ok {
			res, err = f(contextWithUserID(context.Background(), userID), r)
		} else {
			res, err = f(r.Context(), r)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}

		body, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})
}
