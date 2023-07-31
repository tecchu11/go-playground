package handler

import (
	"go-playground/pkg/lib/render"
	"net/http"

	"go.uber.org/zap"
)

type HealthHandler interface {
	GetStatus() http.HandlerFunc
}

type healthHandler struct {
	log *zap.Logger
}

type HealthStatus struct {
	Status string `json:"status"`
}

func NewHealthHandler(log *zap.Logger) HealthHandler {
	return &healthHandler{log: log}
}

func (hh *healthHandler) GetStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Ok(w, &HealthStatus{Status: "OK"})
	}
}
