package datasource_test

import (
	"context"
	"go-playground/cmd/api/internal/datasource"
	"go-playground/cmd/api/internal/datasource/database"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/domain/entity/entitytest"
	"go-playground/pkg/errorx"
	"go-playground/pkg/testhelper"
	"net/http"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserAdaptorFindByID(t *testing.T) {
	type input struct {
		sub string
	}
	type want struct {
		user            entity.User
		isErr           bool
		httpStatusOnErr int
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
			want:  want{isErr: true, httpStatusOnErr: http.StatusNotFound},
		},
	}
	adaptor := datasource.NewUserAdaptor(database.New(db))
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			runInTx(t, func(ctx context.Context) {
				got, err := adaptor.FindBySub(ctx, v.input.sub)
				if v.want.isErr {
					assert.Error(t, err)
					var appErr *errorx.Error
					assert.ErrorAs(t, err, &appErr)
					assert.Equal(t, v.want.httpStatusOnErr, appErr.HTTPStatus())
				} else {
					assert.Equal(t, v.want.user, got)
					assert.NoError(t, err)
				}
			})
		})
	}
}

func TestUserAdaptorCreate(t *testing.T) {
	type input struct {
		user entity.User
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
			input: input{user: entitytest.TestUser(t)},
		},
		"failure duplicate sub": {
			input: input{user: entitytest.TestUser(t, func(u *entity.User) { u.Sub = "80dbb87a-5ce8-4b45-85a0-3b8aec488b7a" })},
			want:  want{isErr: true, httpStatusOnErr: http.StatusBadRequest},
		},
	}
	adaptor := datasource.NewUserAdaptor(database.New(db))
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			runInTx(t, func(ctx context.Context) {
				err := adaptor.Create(ctx, v.input.user)

				if v.want.isErr {
					assert.Error(t, err)
					var appErr *errorx.Error
					assert.ErrorAs(t, err, &appErr)
					assert.Equal(t, v.want.httpStatusOnErr, appErr.HTTPStatus())
				} else {
					got, err := adaptor.FindBySub(ctx, v.input.user.Sub)
					require.NoError(t, err)
					diff := cmp.Diff(got, v.input.user,
						cmpopts.IgnoreFields(got, "CreatedAt", "UpdatedAt"),
					)
					assert.Empty(t, diff)
				}
			})
		})
	}

}
