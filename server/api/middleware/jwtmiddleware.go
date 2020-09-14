package middleware

import (
	"context"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ratedemon/go-rest/config"
)

type JWTMiddleware struct {
	cfg *config.Config
}

type Token struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

func (mware *JWTMiddleware) JWTAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/api/login", "/api/signup"}

		for _, value := range notAuth {
			if value == r.URL.Path {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte("{\"message\": \"Missing Auth Token\"}"))
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte("{\"message\": \"Invalid/Malformed auth token\"}"))
			return
		}

		tokenPart := splitted[1]
		tk := &Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(mware.cfg.JWTSecret), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte("{\"message\": \"Malformed authentication token\"}"))
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			w.Write([]byte("{\"message\": \"Token is not valid\"}"))
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", tk.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func NewJWTMiddleware(cfg *config.Config) *JWTMiddleware {
	return &JWTMiddleware{cfg}
}
