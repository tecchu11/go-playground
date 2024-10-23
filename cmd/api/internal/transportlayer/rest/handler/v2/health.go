package handler

import (
	"context"
	"encoding/json"
	"go-playground/cmd/api/internal/transportlayer/rest/oapi"
	"net/http"
)

// pinger is compatible for [sql.DB]
type pinger interface {
	PingContext(context.Context) error
}

// HealthHandler check service status.
type HealthHandler struct {
	Pinger pinger
}

// HealthCheck is handler for [GET /health].
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	ErrorHandlerFunc(w, r, func(w http.ResponseWriter, r *http.Request) error {
		err := h.Pinger.PingContext(r.Context())
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(oapi.ResponseHealthCheck{Message: "ok"})
	})
}
