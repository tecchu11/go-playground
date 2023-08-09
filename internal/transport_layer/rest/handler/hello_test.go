package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"go-playground/internal/transport_layer/rest/handler"
	"go-playground/internal/transport_layer/rest/middleware"
	"go-playground/internal/transport_layer/rest/model"
	"go-playground/internal/transport_layer/rest/preauth"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestHelloHandler_GetName(t *testing.T) {
	tests := []struct {
		name                string
		inputResponseWriter *httptest.ResponseRecorder
		inputRequest        *http.Request
		expectErr           bool
		expectedCode        int
		expectedBody        model.HelloResponse
		expectedProblem     map[string]string
	}{
		{
			name:                "test GetName returns 200 and expected body",
			inputResponseWriter: httptest.NewRecorder(),
			inputRequest:        httptest.NewRequest("GET", "http://example.com/hello", nil),
			expectErr:           false,
			expectedCode:        200,
			expectedBody:        model.HelloResponse{Message: "Hello tecchu!! You have ADMIN role."},
		},
		{
			name:                "test GetName returns 401 and expected body",
			inputResponseWriter: httptest.NewRecorder(),
			inputRequest:        httptest.NewRequest("GET", "http://example.com/hello", nil),
			expectErr:           true,
			expectedCode:        401,
			expectedProblem:     map[string]string{"title": "Request With No Authentication", "detail": "Request token was not found in your request header"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// mocking middleware.GetAuthenticatedUser
			middleware.CurrentUser = func(ctx context.Context) (*preauth.AuthenticatedUser, error) {
				if test.expectErr {
					return nil, fmt.Errorf("no user")
				}
				return &preauth.AuthenticatedUser{Name: "tecchu", Role: preauth.ADMIN}, nil
			}
			handler.NewHelloHandler(zap.NewExample(), &mockFailure{}).
				GetName().
				ServeHTTP(test.inputResponseWriter, test.inputRequest)

			if test.inputResponseWriter.Code != test.expectedCode {
				t.Errorf("unexpected code(%d) was recived", test.inputResponseWriter.Code)
			}

			if test.expectErr {
				var actualBody map[string]string
				_ = json.Unmarshal(test.inputResponseWriter.Body.Bytes(), &actualBody)
				if !reflect.DeepEqual(actualBody, test.expectedProblem) {
					t.Errorf("unexpected body (%v) was recieved", actualBody)
				}
			} else {
				var actualBody model.HelloResponse
				_ = json.Unmarshal(test.inputResponseWriter.Body.Bytes(), &actualBody)
				if !reflect.DeepEqual(actualBody, test.expectedBody) {
					t.Errorf("unexpected body (%v) was recieved", actualBody)
				}
			}
		})
	}
}
