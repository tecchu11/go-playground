package middleware

import (
	"context"
	"encoding/json"
	"go-playground/pkg/presentation/model"
	"net/http"

	"go.uber.org/zap"
)

type Auth interface {
	Authenticate(key string) (string, error)
}

type AuthKey struct{}

type AuthError struct{}

func (err *AuthError) Error() string {
	return "NO Authenticatted User had been requested."
}

func (a AuthKey) Authenticate(key string) (string, error) {
	if key != "" {
		return "Tetsu", nil
	}
	return "", &AuthError{}
}

func ContextMiddleWareFunc(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "path", r.URL.Path)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func AuthMiddleWareFunc(next http.HandlerFunc, log *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := AuthKey{}.Authenticate(r.Header.Get("Authorization"))
		if err != nil {
			log.Warn("No authenticated request had recieved.", zap.String("path", r.URL.Path), zap.Error(err))
			res := model.NewProblemDetail("You had failed to authenticate", r.URL.Path, http.StatusUnauthorized)
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				log.Error("Failed to write response", zap.Error(err), zap.String("path", r.URL.Path))
			}
		} else {
			ctx := context.WithValue(r.Context(), "user", u)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}
