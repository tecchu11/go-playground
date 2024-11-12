package entity_test

import (
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/pkg/errorx"
	"net/http"
	"testing"

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
		isErr           bool
		httpStatusOnErr int
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
		},
		"failure sub is empty": {
			input: input{
				givenName:     "Walter",
				familyName:    "Sons",
				email:         "Domingo63@example.com",
				emailVerified: true,
			},
			want: want{isErr: true, httpStatusOnErr: http.StatusBadRequest},
		},
		"failure given name is empty": {
			input: input{
				sub:           "f828f57a-b082-4af2-afb4-f3d0fe2c8697",
				familyName:    "Sons",
				email:         "Domingo63@example.com",
				emailVerified: true,
			},
			want: want{isErr: true, httpStatusOnErr: http.StatusBadRequest},
		},
		"failure family name is empty": {
			input: input{
				sub:           "f828f57a-b082-4af2-afb4-f3d0fe2c8697",
				givenName:     "Walter",
				email:         "Domingo63@example.com",
				emailVerified: true,
			},
			want: want{isErr: true, httpStatusOnErr: http.StatusBadRequest},
		},
		"failure email is not valid(not RFC format)": {
			input: input{
				sub:           "f828f57a-b082-4af2-afb4-f3d0fe2c8697",
				givenName:     "Walter",
				familyName:    "Sons",
				email:         "Domingo63..@example.com",
				emailVerified: true,
			},
			want: want{isErr: true, httpStatusOnErr: http.StatusBadRequest},
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			u, err := entity.NewUser(
				v.input.sub,
				v.input.givenName,
				v.input.familyName,
				v.input.email,
				v.input.emailVerified,
			)
			if v.want.isErr {
				assert.Zero(t, u)
				var appErr *errorx.Error
				assert.ErrorAs(t, err, &appErr)
				assert.Equal(t, v.want.httpStatusOnErr, appErr.HTTPStatus())
			} else {
				assert.NotZero(t, u)
				assert.NoError(t, err)
			}
		})
	}

}
