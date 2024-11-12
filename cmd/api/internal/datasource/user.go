package datasource

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-playground/cmd/api/internal/datasource/database"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/domain/repository"
	"go-playground/pkg/errorx"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// UserAdaptor is implementation of [repository.UserRepository]
type UserAdaptor struct {
	queries *database.Queries
}

// NewUserAdaptor create UserAdaptor.
func NewUserAdaptor(q *database.Queries) *UserAdaptor {
	return &UserAdaptor{queries: q}
}

func (a *UserAdaptor) FindBySub(ctx context.Context, sub string) (entity.User, error) {
	defer newrelic.FromContext(ctx).StartSegment("datasource/UserAdaptor/FindByID").End()

	txq := txqFromContext(ctx, a.queries)

	row, err := txq.FindUserBySub(ctx, sub)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, errorx.NewInfo("user is not found", errorx.WithCause(err), errorx.WithStatus(http.StatusNotFound))
		}
		return entity.User{}, errorx.NewError("failed to find user", errorx.WithCause(err))
	}
	uid, err := uuid.FromBytes(row.ID)
	if err != nil {
		return entity.User{}, errorx.NewError(fmt.Sprintf("failed to decode user id(%s)", string(row.ID)), errorx.WithCause(err))
	}
	return entity.User{
		ID:            uid,
		Sub:           row.Sub,
		FamilyName:    row.FamilyName,
		GivenName:     row.GivenName,
		Email:         row.Email,
		EmailVerified: row.EmailVerified,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}, nil
}

// Create creates with given user
func (a *UserAdaptor) Create(ctx context.Context, user entity.User) error {
	defer newrelic.FromContext(ctx).StartSegment("datasource/UserAdaptor/Create").End()

	txq := txqFromContext(ctx, a.queries)

	_, err := txq.CreateUser(ctx, database.CreateUserParams{
		ID:            user.ID[:],
		Sub:           user.Sub,
		GivenName:     user.GivenName,
		FamilyName:    user.FamilyName,
		Email:         user.Email,
		EmailVerified: user.EmailVerified,
	})
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 /* duplicate entry */ {
			return errorx.NewWarn("user is already exist", errorx.WithCause(err), errorx.WithStatus(http.StatusBadRequest))
		}
		return errorx.NewError("failed to create new user", errorx.WithCause(err))
	}
	return nil
}

var _ repository.UserRepository = (*UserAdaptor)(nil)
