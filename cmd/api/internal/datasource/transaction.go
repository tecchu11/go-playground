package datasource

import (
	"context"
	"go-playground/cmd/api/internal/domain/repository"
	"go-playground/pkg/apperr"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type transactionContextKey struct{}

// DBTransactionAdaptor is implementation of repository.TransactionRepository.
type DBTransactionAdaptor struct {
	db *sqlx.DB
}

// NewDBTransactionAdaptor creates pointer of DBTransactionAdaptor.
func NewDBTransactionAdaptor(db *sqlx.DB) *DBTransactionAdaptor {
	return &DBTransactionAdaptor{db}
}

// Do does action in transaction.
func (a *DBTransactionAdaptor) Do(ctx context.Context, action func(context.Context) error) error {
	txx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		return apperr.New("begin db transaction", "unexpected error was happened.", apperr.WithCause(err))
	}
	var done bool
	defer func() {
		if !done {
			if err := txx.Rollback(); err != nil {
				slog.WarnContext(ctx, "failed to rollback", slog.String("error", err.Error()))
			}
		}
	}()
	ctx = context.WithValue(ctx, transactionContextKey{}, txx)
	err = action(ctx)
	if err != nil {
		return err
	}
	done = true
	err = txx.Commit()
	if err != nil {
		return apperr.New("commit to db", "unexpected error was happened.", apperr.WithCause(err))
	}
	return nil
}

var _ repository.TransactionRepository = (*DBTransactionAdaptor)(nil)
