package renderer_test

import (
	"context"
	"encoding/json"
	"go-playground/pkg/renderer"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFailure_Response(t *testing.T) {
	req := renderer.RequestIDFunc(func(_ context.Context) string {
		return "request_id"
	})
	fn := renderer.NewFailure(req).Response

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://example.com/foos", nil)
	title := "this is title"
	detail := "this is detail"

	fn(w, r, http.StatusInternalServerError, title, detail)

	expectedBody := renderer.ProblemDetail{
		Type:      "not:blank",
		Title:     title,
		Detail:    detail,
		Instant:   r.URL.Path,
		RequestID: "request_id",
	}
	expectedCode := 500

	actualCode := w.Code
	var actualBody renderer.ProblemDetail
	_ = json.Unmarshal(w.Body.Bytes(), &actualBody)

	if actualCode != expectedCode {
		t.Errorf("http status code %d was unexpected", actualCode)
	}
	if !reflect.DeepEqual(actualBody, expectedBody) {
		t.Errorf("response body (%v) was unexpected", actualBody)
	}

}
