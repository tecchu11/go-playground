package handler_test

import (
	"encoding/json"
	"go-playground/internal/transport_layer/rest/handler"
	"go-playground/internal/transport_layer/rest/middleware"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHelloHandler_GetName(t *testing.T) {
	tests := map[string]struct {
		requestUser     middleware.AuthUser
		expectErr       bool
		expectedCode    int
		expectedBody    handler.HelloResponse
		expectedErrBody map[string]string
	}{
		"status 200": {
			requestUser:  middleware.AuthUser{Name: "tecchu", Role: middleware.Admin},
			expectErr:    false,
			expectedCode: 200,
			expectedBody: handler.HelloResponse{Message: "Hello tecchu!! You have Admin role."},
		},
		"status 401 when no current user": {
			requestUser:     middleware.NoUser,
			expectErr:       true,
			expectedCode:    401,
			expectedErrBody: map[string]string{"title": "Request With No Authentication", "detail": "Request token was not found in your request header"},
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/hello", nil)
			ctx := v.requestUser.Set(r.Context())
			handler.NewHelloHandler(zap.NewExample(), &mockJSON{}).GetName().ServeHTTP(w, r.WithContext(ctx))

			actualCode := w.Code
			assert.Equal(t, v.expectedCode, actualCode, "status code should be equal")

			if v.expectErr {
				var actualBody map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &actualBody)
				assert.NoError(t, err, "json unmarshal should not be err")
				assert.Equal(t, v.expectedErrBody, actualBody, "error response body should be equal")
				return
			}

			var actualBody handler.HelloResponse
			err := json.Unmarshal(w.Body.Bytes(), &actualBody)
			assert.NoError(t, err, "json unmarshal should not be err")
			assert.Equal(t, v.expectedBody, actualBody, "response body should be equal")

		})
	}
}
