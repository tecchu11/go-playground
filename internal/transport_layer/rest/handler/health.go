package handler

import (
	"go-playground/pkg/renderer"
	"net/http"
)

var status = map[string]string{"status": "ok"}

func StatusHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		renderer.Ok(w, status)
	})
}
