package auth

import (
	"fmt"
	"go-playground/config"
)

var ErrAuthentication = fmt.Errorf("no authenticated user had been requested")

// AuthenticatedUser is represented authenticated user struct.
type AuthenticatedUser struct {
	Name string
	Role Role
}

// AuthenticationManager perform to authenticate user.
type AuthenticationManager interface {
	Authenticate(token string) (*AuthenticatedUser, error)
}

type authenticationManager struct {
	configs []config.AuthConfig
}

func NewAutheticatonManager(configs []config.AuthConfig) AuthenticationManager {
	return &authenticationManager{configs: configs}
}

// Authenticate user with passed token. Ant then, store AuthenticatedUser in context.Context.
func (manager *authenticationManager) Authenticate(token string) (*AuthenticatedUser, error) {
	var ok bool
	var user AuthenticatedUser
	for _, v := range manager.configs {
		if v.Key == token {
			ok = true
			user = AuthenticatedUser{Name: v.Name, Role: RoleFrom(v.RoleStr)}
			break
		}
	}
	if !ok {
		return nil, ErrAuthentication
	}
	return &user, nil
}
