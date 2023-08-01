package middleware_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-playground/pkg/lib/render"
	"go-playground/pkg/presentation/auth"
	"go-playground/pkg/presentation/middleware"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestAuthMiddleWare_Handle(t *testing.T) {
	tests := []struct {
		name                     string
		inputResponseWriter      *httptest.ResponseRecorder
		inputRequest             *http.Request
		inputAuthorizationHeader string
		expectedResponse         string
		expectErr                bool
		expectedErrCode          int
		expectedErrBody          render.ProblemDetail
	}{
		{
			name:                     "test of successful to authenticate and then set user to context",
			inputResponseWriter:      httptest.NewRecorder(),
			inputRequest:             httptest.NewRequest("GET", "http://example.com/anys", nil),
			inputAuthorizationHeader: "valid-token",
			expectedResponse:         "user is tecchu and role is ADMIN",
			expectErr:                false,
		},
		{
			name:                     "test that requests having invalid token will be handled to 401",
			inputResponseWriter:      httptest.NewRecorder(),
			inputRequest:             httptest.NewRequest("GET", "http://example.com/anys", nil),
			inputAuthorizationHeader: "invalid",
			expectErr:                true,
			expectedErrCode:          401,
			expectedErrBody: render.ProblemDetail{
				Type:    "https://github.com/tecchu11/go-playground",
				Title:   "Unauthorized",
				Detail:  "You had failed to authenticate",
				Instant: "/anys",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.inputRequest.Header.Set("Authorization", test.inputAuthorizationHeader)
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				user, _ := middleware.GetAutenticatedUser(r.Context())
				_, _ = fmt.Fprintf(w, "user is %s and role is %s", user.Name, user.Role.String())
			})
			middleware.
				NewAuthMiddleWare(zap.NewExample(), newMockAuthenticationManager()).
				Handle(next).
				ServeHTTP(test.inputResponseWriter, test.inputRequest)

			if !test.expectErr {
				actual := test.inputResponseWriter.Body.String()
				if actual != test.expectedResponse {
					t.Errorf("request context has unexpected user (%v)", actual)
				}
			} else {
				actualCode := test.inputResponseWriter.Code
				var actualBody render.ProblemDetail
				_ = json.Unmarshal(test.inputResponseWriter.Body.Bytes(), &actualBody)

				if actualCode != test.expectedErrCode {
					t.Errorf("unexpected code (%d) was recived", actualCode)
				}
				if !reflect.DeepEqual(actualBody, test.expectedErrBody) {
					t.Errorf("unexpected body (%v) was recieved", actualBody)
				}
			}
		})
	}
}

type mockAuthenticationManager struct{}

func newMockAuthenticationManager() auth.AuthenticationManager {
	return &mockAuthenticationManager{}
}

func (mock *mockAuthenticationManager) Authenticate(token string) (*auth.AuthenticatedUser, error) {
	if token == "valid-token" {
		return &auth.AuthenticatedUser{Name: "tecchu", Role: auth.ADMIN}, nil
	}
	return nil, errors.New("mock")
}
