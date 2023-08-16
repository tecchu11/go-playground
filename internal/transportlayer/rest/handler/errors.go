package handler

import (
	"fmt"
	"go-playground/pkg/renderer"
	"net/http"
)

// NotFoundHandler handle request to resource that does not exist
func NotFoundHandler(rj renderer.JSON) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		detail := fmt.Sprintf("%s resource does not exist", r.URL.Path)
		rj.Failure(w, r, http.StatusNotFound, "Resource Not Found", detail)
	}
}

// MethodNotAllowedHandler handle invalid http method.
func MethodNotAllowedHandler(rj renderer.JSON) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		detail := fmt.Sprintf("Http method %s is not allowed for %s resource", r.Method, r.URL.Path)
		rj.Failure(w, r, http.StatusMethodNotAllowed, "Method Not Allowed", detail)
	}
}
