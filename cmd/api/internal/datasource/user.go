package datasource

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-playground/cmd/api/internal/datasource/database"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/domain/repository"
	"go-playground/pkg/apperr"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// UserAdaptor is implementation of [repository.UserRepository]
type UserAdaptor struct {
	base
}

// NewUserAdaptor create UserAdaptor.
func NewUserAdaptor(db *sqlx.DB) *UserAdaptor {
	return &UserAdaptor{base: base{db: db}}
}

func (a *UserAdaptor) FindBySub(ctx context.Context, sub string) (entity.User, error) {
	defer newrelic.FromContext(ctx).StartSegment("datasource/UserAdaptor/FindByID").End()

	txq := a.queriesFromContext(ctx)

	row, err := txq.FindUserBySub(ctx, sub)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, apperr.New("find user by sub but result set is zero", "user is not found", apperr.WithCause(err), apperr.CodeNotFound)
		}
		return entity.User{}, apperr.New("find user by sub", "failed to find user", apperr.WithCause(err), apperr.CodeNotFound)
	}
	uid, err := uuid.FromBytes(row.ID)
	if err != nil {
		return entity.User{}, apperr.New(fmt.Sprintf("raw user id(%s) to uuid", string(row.ID)), "failed to find user", apperr.WithCause(err))
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

	txq := a.queriesFromContext(ctx)

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
			return apperr.New("create user but user is already exist", "user is already exist", apperr.WithCause(err), apperr.CodeInvalidArgument)
		}
		return apperr.New("create user", "failed to create user", apperr.WithCause(err))
	}
	return nil
}

var _ repository.UserRepository = (*UserAdaptor)(nil)
