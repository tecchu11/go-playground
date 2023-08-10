package renderer_test

import (
	"context"
	"encoding/json"
	"go-playground/pkg/renderer"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOk(t *testing.T) {
	w := httptest.NewRecorder()
	input := map[string]string{"message": "hello"}
	renderer.Ok(w, input)

	expectedCode := 200
	expectedBody := input

	var actualBody map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &actualBody)
	assert.Equal(t, actualBody, expectedBody, "response body should be equal")
	assert.Equal(t, w.Code, expectedCode, "status code should be equal to 200")
}

func TestFailure_Response(t *testing.T) {
	req := renderer.RequestIDFunc(func(_ context.Context) string {
		return "request_id"
	})
	fn := renderer.NewFailure(req).Response
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/foos", nil)
	title := "this is title"
	detail := "this is detail"

	fn(w, r, http.StatusInternalServerError, title, detail)

	expectedBody := renderer.ProblemDetail{Type: "not:blank", Title: title, Detail: detail, Instant: r.URL.Path, RequestID: "request_id"}
	expectedCode := 500

	actualCode := w.Code
	var actualBody renderer.ProblemDetail
	_ = json.Unmarshal(w.Body.Bytes(), &actualBody)

	assert.Equal(t, actualBody, expectedBody, "response body should be equal")
	assert.Equal(t, actualCode, expectedCode, "statuses code should be equal")

}
