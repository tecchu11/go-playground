package datasource_test

import (
	"context"
	"go-playground/cmd/api/internal/datasource"
	"go-playground/cmd/api/internal/datasource/database"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/domain/entity/entitytest"
	"go-playground/pkg/apperr"
	"go-playground/pkg/testhelper"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserAdaptor_FindByID(t *testing.T) {
	type input struct {
		sub string
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
			input: input{sub: "80dbb87a-5ce8-4b45-85a0-3b8aec488b7a"},
			want: want{
				user: entity.User{
					ID:            testhelper.UUIDFromString(t, "01930c3a-e82b-700a-b41a-6f58b5c2b812"),
					Sub:           "80dbb87a-5ce8-4b45-85a0-3b8aec488b7a",
					GivenName:     "Dibbert",
					FamilyName:    "Kozey",
					Email:         "Jonathan74@example.com",
					EmailVerified: true,
					CreatedAt:     time.Date(2024, 11, 8, 14, 40, 33, 0, time.UTC),
					UpdatedAt:     time.Date(2024, 11, 8, 14, 40, 33, 0, time.UTC),
				},
			},
		},
		"failure user not found": {
			input: input{sub: "invalid-sub"},
			want:  want{err: "find user by sub but result set is zero: sql: no rows in result set", errCode: apperr.CodeNotFound},
		},
	}
	adaptor := datasource.NewUserAdaptor(database.New(db))
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runInTx(t, func(ctx context.Context) {
				got, err := adaptor.FindBySub(ctx, tc.input.sub)

				if tc.want.err != "" {
					assert.Zero(t, got)
					assert.EqualError(t, err, tc.want.err)
					assert.True(t, apperr.IsCode(err, tc.want.errCode))
				} else {
					assert.Equal(t, tc.want.user, got)
					assert.NoError(t, err)
				}
			})
		})
	}
}

func TestUserAdaptor_Create(t *testing.T) {
	type input struct {
		user entity.User
	}
	type want struct {
		err     string
		errCode apperr.Code
	}
	tests := map[string]struct {
		input input
		want  want
	}{
		"success": {
			input: input{user: entitytest.TestUser(t)},
		},
		"failure duplicate sub": {
			input: input{user: entitytest.TestUser(t, func(u *entity.User) { u.Sub = "80dbb87a-5ce8-4b45-85a0-3b8aec488b7a" })},
			want:  want{err: "create user but user is already exist: Error 1062 (23000): Duplicate entry '80dbb87a-5ce8-4b45-85a0-3b8aec488b7a' for key 'users.idx_sub'", errCode: apperr.CodeInvalidArgument},
		},
	}
	adaptor := datasource.NewUserAdaptor(database.New(db))
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runInTx(t, func(ctx context.Context) {
				err := adaptor.Create(ctx, tc.input.user)

				if tc.want.err != "" {
					assert.EqualError(t, err, tc.want.err)
					assert.True(t, apperr.IsCode(err, tc.want.errCode))
				} else {
					assert.NoError(t, err)
					got, err := adaptor.FindBySub(ctx, tc.input.user.Sub)
					require.NoError(t, err)
					diff := cmp.Diff(got, tc.input.user, cmpopts.IgnoreFields(entity.User{}, "CreatedAt", "UpdatedAt"))
					assert.Empty(t, diff)
				}
			})
		})
	}

}
