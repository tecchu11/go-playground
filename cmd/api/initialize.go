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

	"go-playground/pkg/nrhttp"
	"go-playground/pkg/nrslog"
	"go-playground/pkg/timex"
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
	myConf := mysql.Config{
		User:         conf.DBUser,
		Passwd:       conf.DBPassword,
		Net:          "tcp",
		Addr:         conf.DBAddr,
		DBName:       conf.DBName,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Loc:          timex.JST(),
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
	// init handler
	listTasks := handler.Listtasks(taskUseCase)
	findTaskByID := handler.FindTaskByID(taskUseCase)
	postTask := handler.PostTask(taskUseCase)
	putTask := handler.PutTask(taskUseCase)
	// init middleware
	middlewares := func(h http.Handler) http.Handler {
		return nrhttp.Middleware(app)(middleware.Recover(h))
	}

	// init router
	mux := http.NewServeMux()
	mux.Handle("GET /health", middlewares(handler.HealthCheck))
	mux.Handle("GET /tasks", middlewares(listTasks))
	mux.Handle("GET /tasks/{id}", middlewares(findTaskByID))
	mux.Handle("POST /tasks", middlewares(postTask))
	mux.Handle("PUT /tasks/{id}", middlewares(putTask))

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
