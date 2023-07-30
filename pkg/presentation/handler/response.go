package handler

import (
	"encoding/json"
	"net/http"
)

func Ok[T any](w http.ResponseWriter, body *T) {
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(body)
}

type ProblemDetail struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Detail  string `json:"detail"`
	Instant string `json:"instant"`
}

func Unauthorized(w http.ResponseWriter, detail string, path string) {
	w.WriteHeader(http.StatusUnauthorized)
	body := &ProblemDetail{
		Type:    "",
		Title:   "Unauthorized",
		Detail:  detail,
		Instant: path,
	}
	_ = json.NewEncoder(w).Encode(body)
}

func NotFound(w http.ResponseWriter, detail string, path string) {
	w.WriteHeader(http.StatusNotFound)
	body := &ProblemDetail{
		Type:    "",
		Title:   "Resource Not Found",
		Detail:  detail,
		Instant: path,
	}
	_ = json.NewEncoder(w).Encode(body)
}

func InternalServerError(w http.ResponseWriter, detail string, path string) {
	w.WriteHeader(http.StatusInternalServerError)
	body := &ProblemDetail{
		Type:    "",
		Title:   "Internal Server Error",
		Detail:  detail,
		Instant: path,
	}
	_ = json.NewEncoder(w).Encode(body)
}
