package handler

import (
	"encoding/json"
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
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(HealthStatus{Status: "OK"}); err != nil {
			hh.log.Error("Failed to write response", zap.Error(err), zap.String("path", r.URL.Path))
		}
	}
}
