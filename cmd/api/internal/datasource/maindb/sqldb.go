package maindb

import (
	"database/sql"
	"go-playground/pkg/env/v2"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/newrelic/go-agent/v3/integrations/nrmysql"
)

// NewQueryDB creates [sql.DB] and [Queries]. Default lookup is [os.LookupEnv].
// Error will be returned if look up failed.
//
// Lookup key is below.
//
//   - DB_USER: user name
//   - DB_PASSWORD: password
//   - DB_ADDRESS: host and port
//   - DB_NAME: database name
func NewQueryDB(lookup func(string) (string, bool)) (*sql.DB, *Queries, error) {
	applier := env.New(lookup)
	conf := mysql.Config{
		User:                 applier.String("DB_USER"),
		Passwd:               applier.String("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 applier.String("DB_ADDRESS"),
		DBName:               applier.String("DB_NAME"),
		MaxAllowedPacket:     64 << 20,
		Timeout:              20 * time.Second,
		ReadTimeout:          20 * time.Second,
		WriteTimeout:         20 * time.Second,
		AllowNativePasswords: true,
		MultiStatements:      true,
		ParseTime:            true,
	}
	if err := applier.Err(); err != nil {
		return nil, nil, err
	}
	db, err := sql.Open("nrmysql", conf.FormatDSN())
	if err != nil {
		return nil, nil, err
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(1 * time.Minute)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db, New(db), nil
}
