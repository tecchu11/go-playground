package middleware_test

import (
	"encoding/json"
	"go-playground/internal/transport_layer/rest/preauth"
	"net/http"

	"github.com/stretchr/testify/mock"
)

// mockFailure mock for renderer.Failure
type mockFailure struct{}

func (m *mockFailure) Response(w http.ResponseWriter, _ *http.Request, code int, title string, detail string) {
	res := map[string]string{"title": title, "detail": detail}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(res)
}

type mockAuthenticationManager struct {
	mock.Mock
}

func (m *mockAuthenticationManager) Authenticate(token string) (*preauth.AuthenticatedUser, error) {
	args := m.Called(token)
	return args.Get(0).(*preauth.AuthenticatedUser), args.Error(1)
}
