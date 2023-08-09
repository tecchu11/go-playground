// render is helper functions for rendering json repsonse with status code.
package render

import (
	"encoding/json"
	"net/http"
)

const (
	contentTypeKey   = "Content-Type"
	contentTypeValue = "application/json; charset=utf-8"
	title            = "https://github.com/tecchu11/go-playground"
)

var (
	statusMap = map[int]string{
		http.StatusUnauthorized:        "Unauthorized",
		http.StatusForbidden:           "Forbidden",
		http.StatusNotFound:            "Resource Not Found",
		http.StatusMethodNotAllowed:    "Method Not Allowed",
		http.StatusInternalServerError: "Internal Server Error",
	}
)

// Ok return 200 and passed body.
func Ok(w http.ResponseWriter, body any) {
	w.Header().Add(contentTypeKey, contentTypeValue)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(body)
}

// ProblemDetail indicates the respose body in accordance with RFC7807.
// Please see detail bellow. https://datatracker.ietf.org/doc/html/rfc7807
type ProblemDetail struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Detail  string `json:"detail"`
	Instant string `json:"instant"`
}

// Unauthorized return 401 and respose body in accordance with ProblemDetail.
func Unauthorized(w http.ResponseWriter, detail string, path string) {
	error(w, detail, path, http.StatusUnauthorized)
}

// Forbidden return 403 and respose body in accordance with ProblemDetail.
func Forbidden(w http.ResponseWriter, detail string, path string) {
	error(w, detail, path, http.StatusForbidden)
}

// NotFound return 404 and respose body in accordance with ProblemDetail.
func NotFound(w http.ResponseWriter, detail string, path string) {
	error(w, detail, path, http.StatusNotFound)
}

// MethodNotAllowed return 405 and respose body in accordance with ProblemDetail.
func MethodNotAllowed(w http.ResponseWriter, detail string, path string) {
	error(w, detail, path, http.StatusMethodNotAllowed)
}

// InternalServerError return 500 and respose body in accordance with ProblemDetail.
func InternalServerError(w http.ResponseWriter, detail string, path string) {
	error(w, detail, path, http.StatusInternalServerError)
}

func error(w http.ResponseWriter, detail string, path string, code int) {
	w.Header().Add(contentTypeKey, contentTypeValue)
	w.WriteHeader(code)
	body := &ProblemDetail{
		Type:    title,
		Title:   statusMap[code],
		Detail:  detail,
		Instant: path,
	}
	_ = json.NewEncoder(w).Encode(body)
}
