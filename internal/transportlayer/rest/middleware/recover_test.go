package middleware_test

import (
	"encoding/json"
	"go-playground/internal/transportlayer/rest/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_RecoverMiddleWare_Handle(t *testing.T) {
	tests := map[string]struct {
		expectErr    bool
		expectedCode int
		expectedBody map[string]string
	}{
		"return 500 and expect body when panic": {
			expectErr:    true,
			expectedCode: 500,
			expectedBody: map[string]string{"title": "Internal Server Error", "detail": "Unexpected error was happened. Please report this error you have checked."},
		},
		"recover middleware nothing to do": {
			expectErr:    false,
			expectedCode: 200,
			expectedBody: map[string]string{"hello": "world"}},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/foos", nil)
			rec := middleware.Recover(zap.NewExample(), &mockJSON{})
			panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if v.expectErr {
					panic("test panic!!")
				}
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(map[string]string{"hello": "world"})
			})
			rec(panicHandler).ServeHTTP(w, r)

			actualCode := w.Code
			var actual map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &actual)

			assert.Equal(t, v.expectedCode, actualCode, "status code should be equal")
			assert.NoError(t, err, "json unmarshal should not be err")
			assert.Equal(t, v.expectedBody, actual, "error response body should be equal")
		})
	}
}
