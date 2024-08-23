package datasource_test

import (
	"context"
	"database/sql"
	"go-playground/cmd/api/internal/datasource"
	"go-playground/pkg/migration"
	"go-playground/pkg/timex"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/go-testfixtures/testfixtures/v3"
	mysqlcontainer "github.com/testcontainers/testcontainers-go/modules/mysql"
)

var db *sql.DB

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	container, shutdown, err := startContainer(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = shutdown(ctx)
	}()
	db = container
	err = db.PingContext(ctx)
	if err != nil {
		panic(err)
	}
	err = migration.Up(ctx, db)
	if err != nil {
		panic(err)
	}
	err = prepareFixture(db)
	if err != nil {
		panic(err)
	}
	m.Run()
}

type shutdownFunc func(context.Context) error

var (
	testImage    = "mysql:8.0.36"
	testDB       = "playground-test"
	testUser     = "test_user"
	testPassword = "test"
)

func startContainer(ctx context.Context) (*sql.DB, shutdownFunc, error) {
	container, err := mysqlcontainer.Run(
		ctx,
		testImage,
		mysqlcontainer.WithDatabase(testDB),
		mysqlcontainer.WithUsername(testUser),
		mysqlcontainer.WithPassword(testPassword),
		mysqlcontainer.WithConfigFile("../../../../testdata/mysql/conf/my.cnf"),
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
	conf := mysql.Config{
		User:                 testUser,
		Passwd:               testPassword,
		Net:                  "tcp",
		Addr:                 addr,
		DBName:               testDB,
		Loc:                  timex.JST(),
		AllowNativePasswords: true,
		ParseTime:            true,
		MultiStatements:      true,
	}
	db, err := sql.Open("mysql", conf.FormatDSN())
	if err != nil {
		_ = container.Terminate(ctx)
		return nil, nil, err
	}
	return db, container.Terminate, nil
}

func prepareFixture(db *sql.DB) error {
	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("mysql"),
		testfixtures.Directory("../../../../testdata/mysql/fixtures"),
	)
	if err != nil {
		return err
	}
	err = fixtures.Load()
	if err != nil {
		return err
	}
	return nil
}

func runInTx(t *testing.T, target func(context.Context)) {
	t.Helper()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("cant start transaction because %v", err)
	}
	defer func() {
		err := tx.Rollback()
		if err != nil {
			t.Fatalf("cant rollback because %v", err)
		}
	}()
	ctx := context.WithValue(context.Background(), datasource.TransactionContextKey{}, tx)
	target(ctx)
}
