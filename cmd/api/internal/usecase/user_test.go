package usecase_test

import (
	"context"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/usecase"
	"go-playground/pkg/errorx"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserUseCase_CreateUser(t *testing.T) {
	type input struct {
		sub, givenName, familyName, email string
		emailVerified                     bool
	}
	type want struct {
		isErr           bool
		httpStatusOnErr int
	}
	type setup func(t *testing.T, input input) *MockUserRepository
	tests := map[string]struct {
		input input
		setup setup
		want  want
	}{
		"success": {
			input: input{
				sub:           "test-sub",
				givenName:     "givenName",
				familyName:    "familyName",
				email:         "email@example.com",
				emailVerified: true,
			},
			setup: func(t *testing.T, input input) *MockUserRepository {
				mck := new(MockUserRepository)
				userMatcher := mock.MatchedBy(func(user entity.User) bool {
					require.Equal(t, input.sub, user.Sub)
					require.Equal(t, input.givenName, user.GivenName)
					require.Equal(t, input.familyName, user.FamilyName)
					require.Equal(t, input.email, user.Email)
					require.Equal(t, input.emailVerified, user.EmailVerified)
					return true
				})
				mck.On("Create", context.Background(), userMatcher).Return(nil)
				return mck
			},
		},
		"failure failed to create user": {
			input: input{
				sub:           "test-sub",
				givenName:     "givenName",
				familyName:    "familyName",
				email:         "email@example.com",
				emailVerified: true,
			},
			setup: func(t *testing.T, input input) *MockUserRepository {
				mck := new(MockUserRepository)
				userMatcher := mock.MatchedBy(func(user entity.User) bool {
					require.Equal(t, input.sub, user.Sub)
					require.Equal(t, input.givenName, user.GivenName)
					require.Equal(t, input.familyName, user.FamilyName)
					require.Equal(t, input.email, user.Email)
					require.Equal(t, input.emailVerified, user.EmailVerified)
					return true
				})
				mck.On("Create", context.Background(), userMatcher).Return(errorx.NewError("error on save user"))
				return mck
			},
			want: want{isErr: true, httpStatusOnErr: http.StatusInternalServerError},
		},
		"failure validation error": {
			setup: func(t *testing.T, input input) *MockUserRepository { return nil },
			want:  want{isErr: true, httpStatusOnErr: http.StatusBadRequest},
		},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			mck := v.setup(t, v.input)
			u := usecase.NewUserUseCase(mck)

			got, err := u.CreateUser(context.Background(), v.input.sub, v.input.givenName, v.input.familyName, v.input.email, v.input.emailVerified)

			if v.want.isErr {
				assert.Zero(t, got)
				var appErr *errorx.Error
				assert.ErrorAs(t, err, &appErr)
				assert.Equal(t, v.want.httpStatusOnErr, appErr.HTTPStatus())
			} else {
				assert.NotZero(t, got)
				assert.NoError(t, err)
			}
		})
	}
}
