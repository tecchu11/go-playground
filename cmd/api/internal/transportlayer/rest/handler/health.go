package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

// Pinger can ping to data source just like [database/sql.DB].
type Pinger interface {
	PingContext(context.Context) error
}

// HealthCheck is handler for health check.
func HealthCheck(pinger Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := pinger.PingContext(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "ping error", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(nil)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"msg": "ok"})
	}
}
