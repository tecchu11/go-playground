package main

import (
	"bufio"
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"go-playground/pkg/migration"
	"go-playground/pkg/timex"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func main() {
	flag.Parse()
	var (
		mode = flag.Arg(0)
		file = flag.Arg(1)
	)

	envs, err := loadDotEnv(file)
	if err != nil {
		panic(err)
	}
	mysql.NewConfig()
	conf := mysql.Config{
		User:      envs["DB_USER"],
		Passwd:    envs["DB_PASSWORD"],
		Net:       "tcp",
		Addr:      envs["DB_HOST"],
		DBName:    envs["DB_NAME"],
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

func loadDotEnv(file string) (map[string]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	envs := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") || !strings.HasPrefix(line, "DB_") {
			continue
		}
		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			return nil, errors.New("unexpected .env format")
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		envs[key] = value
	}
	err = scanner.Err()
	if err != nil {
		return nil, err
	}
	return envs, nil
}
