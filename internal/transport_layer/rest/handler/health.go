package handler

import (
	"go-playground/pkg/renderer"
	"net/http"
)

var status = map[string]string{"status": "ok"}

func StatusHandler(rj renderer.JSON) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rj.Success(w, 200, status)
	}
}
