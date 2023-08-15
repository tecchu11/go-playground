package renderer

import (
	"context"
	"encoding/json"
	"net/http"
)

const (
	contentTypeKey   = "Content-Type"
	contentTypeValue = "application/json; charset=utf-8"
)

type (
	// JSON interface for error response.
	JSON interface {
		Success(w http.ResponseWriter, code int, body any)
		// Failure send to client error response body with status code.
		// Error response body is problemDetail.
		Failure(w http.ResponseWriter, r *http.Request, code int, title string, detail string)
	}
	// RequestIDFunc retrieve request id from context.
	RequestIDFunc func(ctx context.Context) string

	jsonResponse struct {
		requestID RequestIDFunc
	}
	// problemDetail indicates the response body in accordance with RFC7807.
	// Please see detail bellow. https://datatracker.ietf.org/doc/html/rfc7807
	// problemDetail exposed for test purpose.
	problemDetail struct {
		Type      string `json:"type"`
		Title     string `json:"title"`
		Detail    string `json:"detail"`
		Instant   string `json:"instant"`
		RequestID string `json:"request_id"`
	}
)

// NewJSON init JSON.
func NewJSON(requestID func(ctx context.Context) string) JSON {
	return &jsonResponse{requestID}
}

// Failure send to client error response body with status code.
// Error response body is problemDetail.
func (j *jsonResponse) Failure(w http.ResponseWriter, r *http.Request, code int, title string, detail string) {
	reqID := j.requestID(r.Context())
	body := &problemDetail{Type: "not:blank", Title: title, Detail: detail, Instant: r.URL.Path, RequestID: reqID}
	w.Header().Add(contentTypeKey, contentTypeValue)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

// Success send to client code and body. If code is 204 or 205, given body will be ignored.
func (j *jsonResponse) Success(w http.ResponseWriter, code int, body any) {
	w.Header().Add(contentTypeKey, contentTypeValue)
	w.WriteHeader(code)
	if code == http.StatusNoContent || code == http.StatusResetContent {
		return
	}
	_ = json.NewEncoder(w).Encode(body)
}
