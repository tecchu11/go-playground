package handler

import (
	"go-playground/pkg/render"
	"net/http"
)

var status = map[string]string{"status": "ok"}

func StatusHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		render.Ok(w, status)
	})
}
