package renderer

import (
	"context"
	"encoding/json"
	"net/http"
)

type (
	// Failure interface for error response.
	Failure interface {
		// Response send to client error response body with status code.
		// Error response body is ProblemDetail.
		Response(w http.ResponseWriter, r *http.Request, code int, title string, detail string)
	}
	// RequestIDFunc retrieve request id from context.
	RequestIDFunc func(ctx context.Context) string

	failure struct {
		requestID RequestIDFunc
	}
	// ProblemDetail indicates the response body in accordance with RFC7807.
	// Please see detail bellow. https://datatracker.ietf.org/doc/html/rfc7807
	// ProblemDetail exposed for test purpose.
	ProblemDetail struct {
		Type      string `json:"type"`
		Title     string `json:"title"`
		Detail    string `json:"detail"`
		Instant   string `json:"instant"`
		RequestID string `json:"request_id"`
	}
)

// NewFailure init Failure.
func NewFailure(requestID func(ctx context.Context) string) Failure {
	return &failure{requestID}
}

// Response send to client error response body with status code.
// Error response body is ProblemDetail.
func (f *failure) Response(w http.ResponseWriter, r *http.Request, code int, title string, detail string) {
	reqID := f.requestID(r.Context())
	body := &ProblemDetail{Type: "not:blank", Title: title, Detail: detail, Instant: r.URL.Path, RequestID: reqID}
	w.Header().Add(contentTypeKey, contentTypeValue)
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}
