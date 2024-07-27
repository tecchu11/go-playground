package main

import (
	"database/sql"
	"go-playground/cmd/api/internal/config"
	"go-playground/cmd/api/internal/datasource"
	"go-playground/cmd/api/internal/datasource/maindb"
	"go-playground/cmd/api/internal/transportlayer/rest/handler"
	"go-playground/cmd/api/internal/transportlayer/rest/middleware"
	"go-playground/cmd/api/internal/usecase"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/newrelic/go-agent/v3/integrations/nrmysql"
	"github.com/newrelic/go-agent/v3/newrelic"

	"go-playground/pkg/nrmux"
	"go-playground/pkg/nrslog"
	"go-playground/pkg/problemdetails"
)

func Initialize() (*http.Server, error) {
	// set up required
	conf := config.Load()
	app, err := newrelic.NewApplication(newrelic.ConfigFromEnvironment())
	if err != nil {
		return nil, err
	}
	nrHandler, err := nrslog.NewJSONHandler(app, &slog.HandlerOptions{AddSource: true})
	if err != nil {
		return nil, err
	}
	slog.SetDefault(slog.New(nrHandler))
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return nil, err
	}
	myConf := mysql.Config{
		User:         conf.DBUser,
		Passwd:       conf.DBPassword,
		Net:          "tcp",
		Addr:         conf.DBAddr,
		DBName:       conf.DBName,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Loc:          loc,
		ParseTime:    true,
	}
	db, err := sql.Open("nrmysql", myConf.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(1 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	// init datasource
	queries := maindb.New(db)
	transactionAdaptor := datasource.NewDBTransactionAdaptor(db)
	taskAdaptor := datasource.NewTaskAdaptor(queries)
	// useCase
	taskUseCase := usecase.NewTaskUseCase(taskAdaptor, transactionAdaptor)

	// init router
	mux := nrmux.New(app,
		nrmux.WithMarshalJSON404(
			func(r *http.Request) ([]byte, error) {
				return problemdetails.New("Resource Not Found", http.StatusNotFound).JSON(r)
			},
		),
		nrmux.WithMarshalJSON405(
			func(r *http.Request) ([]byte, error) {
				return problemdetails.New("Method Not Allowed", http.StatusMethodNotAllowed).JSON(r)
			},
		),
	)
	mux.Handle("GET /health", middleware.Recover(handler.HealthCheck))
	mux.Handle("GET /tasks/{id}", middleware.Recover(handler.FindTaskByID(taskUseCase)))
	mux.Handle("POST /tasks", middleware.Recover(handler.PostTask(taskUseCase)))
	mux.Handle("PUT /tasks/{id}", middleware.Recover(handler.PutTask(taskUseCase)))

	// inits server
	svr := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	return svr, nil
}
