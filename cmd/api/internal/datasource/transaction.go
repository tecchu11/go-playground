package datasource

import (
	"context"
	"database/sql"
	"go-playground/cmd/api/internal/datasource/database"
	"go-playground/cmd/api/internal/domain/repository"
	"go-playground/pkg/apperr"
	"log/slog"
)

type transactionContextKey struct{}

// txqFromContext configures *sql.tx to given queries if *sql.tx exists in context.
func txqFromContext(ctx context.Context, queries *database.Queries) *database.Queries {
	if tx, ok := ctx.Value(transactionContextKey{}).(*sql.Tx); ok {
		return queries.WithTx(tx)
	}
	return queries
}

// DBTransactionAdaptor is implementation of repository.TransactionRepository.
type DBTransactionAdaptor struct {
	db *sql.DB
}

// NewDBTransactionAdaptor creates pointer of DBTransactionAdaptor.
func NewDBTransactionAdaptor(db *sql.DB) *DBTransactionAdaptor {
	return &DBTransactionAdaptor{db}
}

// Do does action in transaction.
func (a *DBTransactionAdaptor) Do(ctx context.Context, action func(context.Context) error) error {
	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return apperr.New("begin db transaction", "unexpected error was happened.", apperr.WithCause(err))
	}
	var done bool
	defer func() {
		if !done {
			if err := tx.Rollback(); err != nil {
				slog.WarnContext(ctx, "failed to rollback", slog.String("error", err.Error()))
			}
		}
	}()
	ctx = context.WithValue(ctx, transactionContextKey{}, tx)
	err = action(ctx)
	if err != nil {
		return err
	}
	done = true
	err = tx.Commit()
	if err != nil {
		return apperr.New("commit to db", "unexpected error was happened.", apperr.WithCause(err))
	}
	return nil
}

var _ repository.TransactionRepository = (*DBTransactionAdaptor)(nil)
