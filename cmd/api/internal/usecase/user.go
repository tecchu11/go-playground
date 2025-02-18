package usecase

import (
	"context"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// UserUseCase handles user entity.
type UserUseCase struct {
	userRepository repository.UserRepository
}

// NewUserUseCase creates UserUseCase
func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepository: userRepo}
}

func (u *UserUseCase) FindBySub(
	ctx context.Context,
	sub string,
) (entity.User, error) {
	defer newrelic.FromContext(ctx).StartSegment("usecase/UserUseCase/FindMe").End()

	return u.userRepository.FindBySub(ctx, sub)
}

// CreateUser creates new user with given user information.
func (u *UserUseCase) CreateUser(
	ctx context.Context,
	sub string,
	givenName, familyName string,
	email string,
	emailVerified bool,
) (uuid.UUID, error) {
	defer newrelic.FromContext(ctx).StartSegment("usecase/UserUseCase/CreateUser").End()

	user, err := entity.NewUser(sub, givenName, familyName, email, emailVerified)
	if err != nil {
		return uuid.UUID{}, err
	}
	err = u.userRepository.Create(ctx, user)
	if err != nil {
		return uuid.UUID{}, err
	}
	return user.ID, nil
}
