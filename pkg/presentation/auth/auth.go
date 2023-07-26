package auth

import (
	"errors"
	"go-playground/config"
)

// AuthenticationManager perform to authenticate user.
type AuthenticationManager struct {
	Configs []config.AuthConfig
}

// Authenticate user with passed token. Ant then, store AuthenticatedUser in context.Context.
//
// [TODO] Must to verify token.
func (manager *AuthenticationManager) Authenticate(token string) (*AuthenticatedUser, error) {
	var ok bool
	var user AuthenticatedUser
	for _, v := range manager.Configs {
		if v.Key == token {
			ok = true
			user = AuthenticatedUser{Name: v.Name, Role: RoleFrom(v.RoleStr)}
			break
		}
	}
	if !ok {
		return nil, errors.New("no authenticated user had been requested")
	}
	return &user, nil
}
