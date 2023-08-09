package middleware_test

import (
	"encoding/json"
	"errors"
	"go-playground/internal/transport_layer/rest/middleware"
	"go-playground/internal/transport_layer/rest/preauth"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestAuthenticator(t *testing.T) {
	mockManager := &mockAuthenticationManager{}
	tests := map[string]struct {
		token        string
		doMock       *mock.Call
		expectedCode int
		expectedBody map[string]string
	}{
		"sucess to authenticate": {
			token:        "valid-token",
			doMock:       mockManager.On("Authenticate", "valid-token").Return(&preauth.AuthenticatedUser{Name: "tecchu", Role: preauth.ADMIN}, nil),
			expectedCode: 200,
			expectedBody: map[string]string{"name": "tecchu", "role": "ADMIN"},
		},
		"rejected request and then response 401": {
			token:        "invalid-token",
			doMock:       mockManager.On("Authenticate", "invalid-token").Return(&preauth.AuthenticatedUser{}, errors.New("No Authenticated")),
			expectedCode: 401,
			expectedBody: map[string]string{"title": "Request With No Authentication", "detail": "Request token was not found in your request header"},
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/foos", nil)
			r.Header.Set("Authorization", v.token)
			nextFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				user, _ := middleware.CurrentUser(r.Context())
				body := map[string]string{"name": user.Name, "role": user.Role.String()}
				_ = json.NewEncoder(w).Encode(&body)
			})
			_ = v.doMock
			middleware.Authenticator(zap.NewExample(), mockManager, &mockFailure{})(nextFunc).ServeHTTP(w, r)

			actualCode := w.Code
			var actualBody map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &actualBody)

			assert.NoError(t, err, "json unmarshal should not be err")
			assert.Equal(t, v.expectedCode, actualCode, "status code should be equal")
			assert.Equal(t, v.expectedBody, actualBody, "response body should be equal")
		})
	}
}
