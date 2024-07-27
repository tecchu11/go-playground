package migration

import (
	"context"
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed sqls/*sql
var migrations embed.FS

func setup() error {
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("mysql"); err != nil {
		return err
	}
	return nil
}

func Up(ctx context.Context, db *sql.DB) error {
	if err := setup(); err != nil {
		return err
	}
	if err := goose.UpContext(ctx, db, "sqls"); err != nil {
		return err
	}
	return nil
}

func Down(ctx context.Context, db *sql.DB) error {
	if err := setup(); err != nil {
		return err
	}
	if err := goose.DownContext(ctx, db, "sqls"); err != nil {
		return err
	}
	return nil
}
