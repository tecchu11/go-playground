package preauth_test

import (
	"go-playground/configs"
	"go-playground/internal/transport_layer/rest/preauth"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticationManager_Authenticate(t *testing.T) {
	confs := []configs.AuthConfig{
		{Name: "test-user-1", RoleStr: "ADMIN", Key: "test-api-key-1"},
		{Name: "test-user-2", RoleStr: "USER", Key: "test-api-key-2"},
	}
	manager := preauth.NewAuthenticationManager(confs)
	cases := []struct {
		name         string
		token        string
		expectedUser *preauth.AuthenticatedUser
		expectErr    bool
	}{
		{name: "case of successful to authentication", token: "test-api-key-2", expectedUser: &preauth.AuthenticatedUser{Name: "test-user-2", Role: preauth.USER}},
		{name: "case of failure to authentication", token: "invalid api key", expectErr: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actualUser, actualErr := manager.Authenticate(c.token)
			assert.Equal(t, c.expectedUser, actualUser)
			if c.expectErr {
				assert.Error(t, actualErr, "actual err must not be nil")
			}
		})
	}
}
