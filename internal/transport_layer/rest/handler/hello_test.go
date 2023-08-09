package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"go-playground/internal/transport_layer/rest/handler"
	"go-playground/internal/transport_layer/rest/middleware"
	"go-playground/internal/transport_layer/rest/model"
	"go-playground/internal/transport_layer/rest/preauth"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestHelloHandler_GetName(t *testing.T) {
	tests := map[string]struct {
		expectErr       bool
		expectedCode    int
		expectedBody    model.HelloResponse
		expectedErrBody map[string]string
	}{
		"status 200": {
			expectErr:    false,
			expectedCode: 200,
			expectedBody: model.HelloResponse{Message: "Hello tecchu!! You have ADMIN role."},
		},
		"status 401 when no current user": {
			expectErr:       true,
			expectedCode:    401,
			expectedErrBody: map[string]string{"title": "Request With No Authentication", "detail": "Request token was not found in your request header"},
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			// mocking middleware.CurrentUser
			middleware.CurrentUser = func(ctx context.Context) (*preauth.AuthenticatedUser, error) {
				if v.expectErr {
					return nil, fmt.Errorf("no user")
				}
				return &preauth.AuthenticatedUser{Name: "tecchu", Role: preauth.ADMIN}, nil
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/hello", nil)
			handler.NewHelloHandler(zap.NewExample(), &mockFailure{}).GetName().ServeHTTP(w, r)

			actualCode := w.Code
			assert.Equal(t, v.expectedCode, actualCode, "status code should be equal")

			if v.expectErr {
				var actualBody map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &actualBody)
				assert.NoError(t, err, "json unmarshal should not be err")
				assert.Equal(t, v.expectedErrBody, actualBody, "error response body should be equal")
				return
			}

			var actualBody model.HelloResponse
			err := json.Unmarshal(w.Body.Bytes(), &actualBody)
			assert.NoError(t, err, "json unmarshal should not be err")
			assert.Equal(t, v.expectedBody, actualBody, "response body should be equal")

		})
	}
}
