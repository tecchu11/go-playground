package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"go-playground/pkg/lib/render"
	"go-playground/pkg/presentation/auth"
	"go-playground/pkg/presentation/handler"
	"go-playground/pkg/presentation/model"
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
		inputUser           *auth.AuthenticatedUser
		expectedCode        int
		expectErr           bool
		expectedBody        model.HelloResponse
		expectedProblem     render.ProblemDetail
	}{
		{
			name:                "test GetName returns 200 and expected body",
			inputResponseWriter: httptest.NewRecorder(),
			inputRequest:        httptest.NewRequest("GET", "http://example.com/hello", nil),
			inputUser:           &auth.AuthenticatedUser{Name: "tecchu", Role: auth.ADMIN},
			expectedCode:        200,
			expectedBody:        model.HelloResponse{Message: "Hello tecchu!! You have ADMIN role."},
		},
		{
			name:                "test GetName returns 401 and expected body",
			inputResponseWriter: httptest.NewRecorder(),
			inputRequest:        httptest.NewRequest("GET", "http://example.com/hello", nil),
			expectedCode:        401,
			expectErr:           true,
			expectedProblem: render.ProblemDetail{
				Type:    "",
				Title:   "Unauthorized",
				Detail:  "No token was found for your request",
				Instant: "/hello",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler.
				NewHelloHandler(zap.NewExample(), &mockCtxManager{test.expectErr}).
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

type mockCtxManager struct {
	isErr bool
}

func (m *mockCtxManager) Get(ctx context.Context) (*auth.AuthenticatedUser, error) {
	if m.isErr {
		return nil, fmt.Errorf("no user")
	}
	return &auth.AuthenticatedUser{Name: "tecchu", Role: auth.ADMIN}, nil
}
