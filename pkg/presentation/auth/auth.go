package auth

import (
	"errors"
)

// AuthenticationManager perform to authenticate user.
type AuthenticationManager struct {
}

// Authenticate user with passed token. Ant then, store AuthenticatedUser in context.Context.
//
// [TODO] Must to verify token.
func (manager *AuthenticationManager) Authenticate(token string) (*AuthenticatedUser, error) {
	if len(token) == 0 {
		return nil, errors.New("no authenticated user had been requested")
	}
	user := AuthenticatedUser{Name: "tecchu", Role: ADMIN}
	return &user, nil
}
