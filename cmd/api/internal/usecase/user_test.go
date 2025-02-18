package usecase_test

import (
	"context"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/usecase"
	"go-playground/pkg/apperr"
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
		err     string
		errCode apperr.Code
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
				mck.On("Create", context.Background(), userMatcher).Return(apperr.New("internal server error", "failed to update task"))
				return mck
			},
			want: want{err: "internal server error", errCode: apperr.CodeInternal},
		},
		"failure validation error": {
			setup: func(t *testing.T, input input) *MockUserRepository { return nil },
			want:  want{err: "validate user entity: Email: cannot be blank; FamilyName: cannot be blank; GivenName: cannot be blank; Sub: cannot be blank.", errCode: apperr.CodeInvalidArgument},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mck := tc.setup(t, tc.input)
			u := usecase.NewUserUseCase(mck)

			got, err := u.CreateUser(context.Background(), tc.input.sub, tc.input.givenName, tc.input.familyName, tc.input.email, tc.input.emailVerified)

			if tc.want.err != "" {
				assert.Zero(t, got)
				assert.EqualError(t, err, tc.want.err)
				assert.True(t, apperr.IsCode(err, tc.want.errCode))
			} else {
				assert.NotZero(t, got)
				assert.NoError(t, err)
			}
		})
	}
}

func TestFindMe(t *testing.T) {
	type input struct {
		ctx context.Context
		sub string
	}
	in := input{ctx: context.Background(), sub: "0195195a-2958-7ccd-b39d-c7cbe04128b1"}
	mck := new(MockUserRepository)
	uc := usecase.NewUserUseCase(mck)
	mck.On("FindBySub", in.ctx, in.sub).
		Return(entity.User{}, nil)

	_, err := uc.FindBySub(in.ctx, in.sub)

	assert.NoError(t, err)
}
