package preauth_test

import (
	"go-playground/configs"
	preauth2 "go-playground/internal/interactor/rest/preauth"
	"reflect"
	"testing"
)

func TestAuthenticationManager_Authenticate(t *testing.T) {
	configs := []configs.AuthConfig{
		{
			Name:    "test-user-1",
			RoleStr: "ADMIN",
			Key:     "test-api-key-1",
		},
		{
			Name:    "test-user-2",
			RoleStr: "USER",
			Key:     "test-api-key-2",
		},
	}
	manager := preauth2.NewAutheticatonManager(configs)
	cases := []struct {
		name         string
		token        string
		expectedUser *preauth2.AuthenticatedUser
		expectErr    bool
	}{
		{
			name:         "case of successful to authentication",
			token:        "test-api-key-2",
			expectedUser: &preauth2.AuthenticatedUser{Name: "test-user-2", Role: preauth2.USER},
		},
		{
			name:      "case of failuer to authentication",
			token:     "invalid api key",
			expectErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actualUser, actualErr := manager.Authenticate(c.token)
			if !reflect.DeepEqual(actualUser, c.expectedUser) {
				t.Errorf("Unmatched user. actualUser is %v but expected is %v", actualUser, c.expectedUser)
			}
			if c.expectErr && actualErr == nil {
				t.Errorf(" expected error but actula error is nil")
			}

		})
	}
}
