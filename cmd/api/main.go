package main

import (
	"context"
	"errors"
	"go-playground/cmd/api/internal/datasource"
	"go-playground/cmd/api/internal/datasource/maindb"
	"go-playground/cmd/api/internal/transportlayer/rest/handler"
	"go-playground/cmd/api/internal/transportlayer/rest/middleware"
	"go-playground/cmd/api/internal/usecase"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/tecchu11/nrgo-std/nrhttp"
	"github.com/tecchu11/nrgo-std/nrslog"
)

func main() {
	svr, err := setup()
	if err != nil {
		panic(err)
	}
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		slog.Info("We received an interrupt signal,so attempt to shutdown with gracefully")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := svr.Shutdown(ctx); err != nil {
			slog.Error("Failed to terminate server with gracefully. So force terminating ...", slog.String("error", err.Error()))
		}
		close(idleConnsClosed)
	}()

	slog.Info("Server starting ---(ﾟ∀ﾟ)---!!!")
	if err := svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Failed to start up server", slog.String("error", err.Error()))
		panic(err)
	}
	<-idleConnsClosed
	slog.Info("Bye!!")
}

func setup() (*http.Server, error) {
	app, err := newrelic.NewApplication(newrelic.ConfigFromEnvironment())
	if err != nil {
		return nil, err
	}
	slog.SetDefault(slog.New(
		nrslog.NewHandler(
			app,
			nrslog.WithHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})),
		),
	))

	db, queries, err := maindb.NewQueryDB(os.LookupEnv)
	if err != nil {
		return nil, err
	}
	// init datasource
	transactionAdaptor := datasource.NewDBTransactionAdaptor(db)
	taskAdaptor := datasource.NewTaskAdaptor(queries)
	// useCase
	taskUseCase := usecase.NewTaskUseCase(taskAdaptor, transactionAdaptor)
	// init handler
	health := handler.HealthCheck(db)
	listTasks := handler.ListTasks(taskUseCase)
	findTaskByID := handler.FindTaskByID(taskUseCase)
	postTask := handler.PostTask(taskUseCase)
	putTask := handler.PutTask(taskUseCase)
	// init middleware
	middlewares := func(h http.Handler) http.Handler {
		return nrhttp.Middleware(app)(middleware.Recover(h))
	}

	// init router
	mux := http.NewServeMux()
	mux.Handle("GET /health", middlewares(health))
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
