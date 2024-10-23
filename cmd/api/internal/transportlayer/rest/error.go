package rest

import (
	"encoding/json"
	"go-playground/cmd/api/internal/transportlayer/rest/oapi"
	"net/http"
)

// ErrBody is common error response body.
type ErrBody struct {
	Msg string `json:"message"`
}

// Err writes error response.
func Err(
	w http.ResponseWriter,
	msg string,
	sts int,
) {
	w.WriteHeader(sts)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(oapi.Error{Message: msg})
}
