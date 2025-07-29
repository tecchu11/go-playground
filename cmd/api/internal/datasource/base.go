package datasource

import (
	"context"
	"go-playground/cmd/api/internal/datasource/database"

	"github.com/jmoiron/sqlx"
)

type base struct {
	db *sqlx.DB
}

// queriesFromContext retrieves [database.Queries] configured transaction(if any) from [context.Context].
func (b *base) queriesFromContext(ctx context.Context) database.Queries {
	txx, ok := ctx.Value(transactionContextKey{}).(*sqlx.Tx)
	if !ok {
		return *database.New(b.db)
	}
	return *database.New(txx)
}

// extFromContext retrieves [sqlx.Ext] configured transaction(if any) from [context.Context].
func (b *base) extFromContext(ctx context.Context) sqlx.Ext {
	txx, ok := ctx.Value(transactionContextKey{}).(*sqlx.Tx)
	if !ok {
		return b.db
	}
	return txx
}
