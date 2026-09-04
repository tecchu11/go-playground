// Package innodblock_test measures, against a real MySQL 8.0 instance, which
// InnoDB locks a DELETE takes under REPEATABLE READ.
//
// The question it answers: does
//
//	DELETE FROM job_attachment WHERE (job_id, ord) IN ((1,1),(1,2))
//
// take gap locks, or only record locks? Gap locks are what turn the
// delete-then-insert pattern into a deadlock between two transactions that
// touch unrelated job rows.
package innodblock_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/testcontainers/testcontainers-go"
	mysqlcontainer "github.com/testcontainers/testcontainers-go/modules/mysql"
)

var db *sql.DB

var (
	testImage    = "mysql:8.0.36"
	testDB       = "playground-test"
	testUser     = "test_user"
	testPassword = "test"
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	conn, terminate, err := startContainer(ctx)
	if err != nil {
		panic(err)
	}
	db = conn
	code := 1
	func() {
		defer func() {
			_ = terminate(context.Background())
		}()
		if err := db.PingContext(ctx); err != nil {
			panic(err)
		}
		code = m.Run()
	}()
	os.Exit(code)
}

type shutdownFunc func(context.Context, ...testcontainers.TerminateOption) error

func startContainer(ctx context.Context) (*sql.DB, shutdownFunc, error) {
	container, err := mysqlcontainer.Run(
		ctx,
		testImage,
		mysqlcontainer.WithDatabase(testDB),
		mysqlcontainer.WithUsername(testUser),
		mysqlcontainer.WithPassword(testPassword),
		mysqlcontainer.WithConfigFile("../../testdata/mysql/conf/my.cnf"),
	)
	if err != nil {
		return nil, nil, err
	}
	if err := container.Start(ctx); err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, err
	}
	addr, err := container.Endpoint(ctx, "")
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, err
	}
	// performance_schema.data_locks needs the PROCESS privilege, which the
	// container's application user does not have.
	conf := mysql.Config{
		User:                 "root",
		Passwd:               testPassword,
		Net:                  "tcp",
		Addr:                 addr,
		DBName:               testDB,
		AllowNativePasswords: true,
		ParseTime:            true,
		MultiStatements:      true,
	}
	conn, err := sql.Open("mysql", conf.FormatDSN())
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, err
	}
	return conn, container.Terminate, nil
}
