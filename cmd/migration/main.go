package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"go-playground/pkg/env/v2"
	"go-playground/pkg/migration"
	"go-playground/pkg/timex"

	"github.com/go-sql-driver/mysql"
)

func main() {
	flag.Parse()
	mode := flag.Arg(0)
	if mode != "up" && mode != "down" {
		panic(`mode must be "up" or "down"`)
	}
	applier := env.New(nil)
	conf := mysql.Config{
		User:      applier.String("DB_USER"),
		Passwd:    applier.String("DB_PASSWORD"),
		Net:       "tcp",
		Addr:      applier.String("DB_ADDRESS"),
		DBName:    applier.String("DB_NAME"),
		Loc:       timex.JST(),
		ParseTime: true,
	}
	db, err := sql.Open("mysql", conf.FormatDSN())
	if err != nil {
		panic(err)
	}
	switch mode {
	case "up":
		err := migration.Up(context.Background(), db)
		if err != nil {
			panic(err)
		}
	case "down":
		err := migration.Down(context.Background(), db)
		if err != nil {
			panic(err)
		}
	default:
		err := fmt.Errorf(`mode must be "up" or "down"`)
		panic(err)
	}
}
