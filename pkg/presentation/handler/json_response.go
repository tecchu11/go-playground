package handler

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, status int, body any) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(body)
}
