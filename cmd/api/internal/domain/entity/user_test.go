package entity_test

import (
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/pkg/apperr"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	type input struct {
		sub                   string
		givenName, familyName string
		email                 string
		emailVerified         bool
	}
	type want struct {
		user    entity.User
		err     string
		errCode apperr.Code
	}
	tests := map[string]struct {
		input input
		want  want
	}{
		"success": {
			input: input{
				sub:           "f828f57a-b082-4af2-afb4-f3d0fe2c8697",
				givenName:     "Walter",
				familyName:    "Sons",
				email:         "Domingo63@example.com",
				emailVerified: true,
			},
			want: want{
				user: entity.User{
					Sub:           "f828f57a-b082-4af2-afb4-f3d0fe2c8697",
					GivenName:     "Walter",
					FamilyName:    "Sons",
					Email:         "Domingo63@example.com",
					EmailVerified: true,
				},
			},
		},
		"failure sub is empty": {
			input: input{
				givenName:     "Walter",
				familyName:    "Sons",
				email:         "Domingo63@example.com",
				emailVerified: true,
			},
			want: want{

				err:     "validate user entity: Sub: cannot be blank.",
				errCode: apperr.CodeInvalidArgument,
			},
		},
		"failure given name is empty": {
			input: input{
				sub:           "f828f57a-b082-4af2-afb4-f3d0fe2c8697",
				familyName:    "Sons",
				email:         "Domingo63@example.com",
				emailVerified: true,
			},
			want: want{
				err:     "validate user entity: GivenName: cannot be blank.",
				errCode: apperr.CodeInvalidArgument,
			},
		},
		"failure family name is empty": {
			input: input{
				sub:           "f828f57a-b082-4af2-afb4-f3d0fe2c8697",
				givenName:     "Walter",
				email:         "Domingo63@example.com",
				emailVerified: true,
			},
			want: want{
				err:     "validate user entity: FamilyName: cannot be blank.",
				errCode: apperr.CodeInvalidArgument,
			},
		},
		"failure email is not valid(not RFC format)": {
			input: input{
				sub:           "f828f57a-b082-4af2-afb4-f3d0fe2c8697",
				givenName:     "Walter",
				familyName:    "Sons",
				email:         "Domingo63..@example.com",
				emailVerified: true,
			},
			want: want{
				err:     "validate user entity: Email: must be a valid email address.",
				errCode: apperr.CodeInvalidArgument,
			},
		},
		"failure email is empty": {
			input: input{
				sub:           "f828f57a-b082-4af2-afb4-f3d0fe2c8697",
				givenName:     "Walter",
				familyName:    "Sons",
				emailVerified: true,
			},
			want: want{
				err:     "validate user entity: Email: cannot be blank.",
				errCode: apperr.CodeInvalidArgument,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			u, err := entity.NewUser(
				tc.input.sub,
				tc.input.givenName,
				tc.input.familyName,
				tc.input.email,
				tc.input.emailVerified,
			)
			if tc.want.err != "" {
				assert.Zero(t, u)
				assert.Equal(t, tc.want.err, err.Error())
				assert.True(t, apperr.IsCode(err, tc.want.errCode))
			} else {
				diff := cmp.Diff(tc.want.user, u, cmpopts.IgnoreFields(entity.User{}, "ID", "CreatedAt", "UpdatedAt"))
				assert.Empty(t, diff)
				assert.NoError(t, err)
			}
		})
	}

}
