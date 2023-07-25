package middleware

import (
	"context"
	"net/http"
)

type contextMiddleWare struct {
}

func NewContextMiddleWare() MiddleWare {
	return &contextMiddleWare{}
}

func (mid *contextMiddleWare) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "path", r.URL.Path)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
