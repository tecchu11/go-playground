package preauth

import (
	"errors"
	"go-playground/config"
)

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

// NewAuthenticationManager is factory method for AuthenticationManager.
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
		return nil, errors.New("no authenticated user had been requested")
	}
	return &user, nil
}

// PreAuthenticatedList mangess pre authenticated AuthenticatedUser list.
type PreAuthenticatedMap map[string]AuthenticatedUser

// AuthenticateWith authenticate user requested token by PreAuthenticatedList
func (preAuthenticatedMap PreAuthenticatedMap) AuthenticateWith(token string) (*AuthenticatedUser, error) {
	user, ok := preAuthenticatedMap[token]
	if !ok {
		return nil, errors.New("no authenticated user had been requested")
	}
	return &user, nil
}
