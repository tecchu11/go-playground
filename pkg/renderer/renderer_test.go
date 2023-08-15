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

func TestJsonResponse_Success(t *testing.T) {
	tests := map[string]struct {
		inW         *httptest.ResponseRecorder
		inCode      int
		inBody      map[string]string
		exceptEmpty bool
		wantCode    int
		wantBody    map[string]string
	}{
		"success with status code 200": {inW: httptest.NewRecorder(), inCode: 200, inBody: map[string]string{"message": "hello"}, wantCode: 200, wantBody: map[string]string{"message": "hello"}},
		"success with status code 204": {inW: httptest.NewRecorder(), inCode: 204, inBody: map[string]string{"message": "hello"}, exceptEmpty: true, wantCode: 204},
		"success with status code 205": {inW: httptest.NewRecorder(), inCode: 205, inBody: map[string]string{"message": "hello"}, exceptEmpty: true, wantCode: 205},
	}
	fn := func(ctx context.Context) string {
		return ""
	}
	success := renderer.NewJSON(fn).Success

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			success(v.inW, v.inCode, v.inBody)

			gotCode := v.inW.Code
			assert.Equal(t, v.wantCode, gotCode)
			if v.exceptEmpty {
				gotBody := v.inW.Body.Bytes()
				assert.Emptyf(t, gotBody, "")
				return
			}
			var gotBody map[string]string
			err := json.NewDecoder(v.inW.Body).Decode(&gotBody)
			assert.NoError(t, err)
			assert.Equal(t, v.inBody, gotBody)
		})
	}
}

func TestJsonResponse_Failure(t *testing.T) {
	req := renderer.RequestIDFunc(func(_ context.Context) string {
		return "request_id"
	})
	fn := renderer.NewJSON(req).Failure
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
