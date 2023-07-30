package auth_test

import (
	"errors"
	"go-playground/config"
	"go-playground/pkg/presentation/auth"
	"reflect"
	"testing"
)

func TestAuthenticationManager_Authenticate(t *testing.T) {
	configs := []config.AuthConfig{
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
	manager := auth.NewAutheticatonManager(configs)
	cases := []struct {
		name         string
		token        string
		expectedUser *auth.AuthenticatedUser
		expectedErr  error
	}{
		{
			name:         "case of successful to authentication",
			token:        "test-api-key-2",
			expectedUser: &auth.AuthenticatedUser{Name: "test-user-2", Role: auth.USER},
		},
		{
			name:        "case of failuer to authentication",
			token:       "invalid api key",
			expectedErr: auth.ErrAuthentication,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actualUser, actualErr := manager.Authenticate(c.token)
			if !reflect.DeepEqual(actualUser, c.expectedUser) {
				t.Errorf("Unmatched user. actualUser is %v but expected is %v", actualUser, c.expectedUser)
			}
			if !errors.Is(actualErr, c.expectedErr) {
				t.Errorf("Unmatched error. actualUser is (%v) but expected is (%v)", actualErr, c.expectedErr)
			}

		})
	}
}
