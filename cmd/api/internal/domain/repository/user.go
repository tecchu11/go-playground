package repository

import (
	"context"
	"go-playground/cmd/api/internal/domain/entity"
)

// UserRepository manipulates user datastore.
type UserRepository interface {
	// Create creates new user with given entity of user.
	Create(context.Context, entity.User) error
}