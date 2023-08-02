package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"go-playground/pkg/lib/render"
	"go-playground/pkg/presentation/handler"
	"go-playground/pkg/presentation/middleware"
	"go-playground/pkg/presentation/model"
	"go-playground/pkg/presentation/preauth"
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
		expectedProblem     render.ProblemDetail
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
			expectedProblem: render.ProblemDetail{
				Type:    "https://github.com/tecchu11/go-playground",
				Title:   "Unauthorized",
				Detail:  "No token was found for your request",
				Instant: "/hello",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// mocking middleware.GetAuthenticatedUser
			middleware.GetAutenticatedUser = func(ctx context.Context) (*preauth.AuthenticatedUser, error) {
				if test.expectErr {
					return nil, fmt.Errorf("no user")
				}
				return &preauth.AuthenticatedUser{Name: "tecchu", Role: preauth.ADMIN}, nil
			}
			handler.
				NewHelloHandler(zap.NewExample()).
				GetName().
				ServeHTTP(test.inputResponseWriter, test.inputRequest)

			if test.inputResponseWriter.Code != test.expectedCode {
				t.Errorf("unexpected code(%d) was recived", test.inputResponseWriter.Code)
			}

			if test.expectErr {
				var actualBody render.ProblemDetail
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
