package maindb

import (
	"database/sql"
	"go-playground/pkg/env"
	"go-playground/pkg/timex"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/newrelic/go-agent/v3/integrations/nrmysql"
)

// NewDB creates [sql.DB]. Default lookup is [os.LookupEnv].
// Error will be returned if look up failed.
//
// Lookup key is below.
//
//   - DB_USER: user name
//   - DB_PASSWORD: password
//   - DB_ADDRESS: host and port
//   - DB_NAME: database name
func NewDB(lookup func(string) (string, bool)) (*sql.DB, error) {
	var err error
	conf := mysql.Config{
		User:                 env.ApplyString(&err, "DB_USER", lookup),
		Passwd:               env.ApplyString(&err, "DB_PASSWORD", lookup),
		Net:                  "tcp",
		Addr:                 env.ApplyString(&err, "DB_ADDRESS", lookup),
		DBName:               env.ApplyString(&err, "DB_NAME", lookup),
		Loc:                  timex.JST(),
		MaxAllowedPacket:     64 << 20,
		Timeout:              20 * time.Second,
		ReadTimeout:          20 * time.Second,
		WriteTimeout:         20 * time.Second,
		AllowNativePasswords: true,
		MultiStatements:      true,
		ParseTime:            true,
	}
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("nrmysql", conf.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(1 * time.Minute)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}
