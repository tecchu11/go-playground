package handler

import (
	"fmt"
	"go-playground/pkg/lib/render"
	"net/http"
)

// NotFoundHandler handle request to resource that does not exist
func NotFoundHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.NotFound(w, "A request for a resource that does not exist.", r.URL.Path)
	})
}

// MethodNotAllowedHandler handle invalid http method.
func MethodNotAllowedHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.MethodNotAllowed(w, fmt.Sprintf("Http method %s is not allowed for this resource.", r.Method), r.URL.Path)
	})
}
