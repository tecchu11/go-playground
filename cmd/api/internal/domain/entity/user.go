package entity

import (
	"go-playground/pkg/apperr"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

// User is entity of user.
type User struct {
	ID                    uuid.UUID
	Sub                   string
	GivenName, FamilyName string
	Email                 string
	EmailVerified         bool
	CreatedAt, UpdatedAt  time.Time
}

// NewUser creates new user.
func NewUser(
	sub string,
	givenName, familyName string,
	email string,
	emailVerified bool,
) (User, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return User{}, apperr.New("uuid new v7 for user id", "Failed to create new user", apperr.WithCause(err))
	}
	now := time.Now()
	user := User{
		ID:            uid,
		Sub:           sub,
		GivenName:     givenName,
		FamilyName:    familyName,
		Email:         email,
		EmailVerified: emailVerified,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	err = user.validate()
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// validate validates user entity.
func (u User) validate() error {
	err := validation.ValidateStruct(
		&u,
		validation.Field(&u.Sub, validation.Required),
		validation.Field(&u.GivenName, validation.Required),
		validation.Field(&u.FamilyName, validation.Required),
		validation.Field(&u.Email, validation.Required, is.Email),
	)
	if err != nil {
		return apperr.New("validate user entity", err.Error(), apperr.WithCause(err), apperr.CodeInvalidArgument)
	}

	return nil
}
