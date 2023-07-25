package auth

import (
	"context"
	"errors"
)

const authenticatedUserKey = "authUser"

// AuthenticatedUser is represented authenticated user struct.
type AuthenticatedUser struct {
	Name string
	Role Role
}

// GetAuthUser retrive AuthenticatedUser from context.Context
func GetAuthUser(ctx context.Context) (*AuthenticatedUser, error) {
	u, ok := ctx.Value(authenticatedUserKey).(*AuthenticatedUser)
	if !ok {
		return nil, errors.New("user does not exist context")
	}
	return u, nil
}

// SetContext store AuthenticatedUser in context.Context
func (user *AuthenticatedUser) SetContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, authenticatedUserKey, user)
}
