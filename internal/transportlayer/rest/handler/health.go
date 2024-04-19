package handler

import (
	"encoding/json"
	"net/http"
)

// HealthCheck is handler for health check.
var HealthCheck = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	res := map[string]string{"status": "ok"}
	json.NewEncoder(w).Encode(res)
})
