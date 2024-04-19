package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"

	"go-playground/internal/transportlayer/rest/handler"
	"go-playground/internal/transportlayer/rest/middleware"
	"go-playground/pkg/nrslog"
	"go-playground/pkg/problemdetails"
	"go-playground/pkg/router"
)

func Initialize() (*http.Server, error) {
	// set up required
	env, ok := os.LookupEnv("APP_ENV")
	if !ok {
		return nil, errors.New("environment variables APP_ENV is missing")
	}
	conf, err := LoadConfig(env)
	if err != nil {
		return nil, err
	}
	app, err := newrelic.NewApplication(newrelic.ConfigFromEnvironment())
	if err != nil {
		return nil, err
	}
	nrHandler, err := nrslog.NewJSONHandler(app, &slog.HandlerOptions{AddSource: true})
	if err != nil {
		return nil, err
	}
	slog.SetDefault(slog.New(nrHandler))

	// inits router
	newrelicMiddleware := middleware.NewrelicTxn(app)
	recoverMiddleware := middleware.Recover
	compositeMiddleware := func(next http.Handler) http.Handler {
		return newrelicMiddleware(recoverMiddleware(next))
	}
	// init router
	mux := router.New(
		router.With404Func(func(r *http.Request) ([]byte, error) {
			return problemdetails.New("Resource Not Found", http.StatusNotFound).JSON(r)
		}),
		router.With405Func(func(r *http.Request) ([]byte, error) {
			return problemdetails.New("Method Not Allowed", http.StatusMethodNotAllowed).JSON(r)
		}),
		router.WithMiddleware(compositeMiddleware),
	)
	mux.Handle("GET /health", handler.HealthCheck)
	mux.Handle("GET /reply/{name}", handler.ReplyHandler)

	// inits server
	svr := &http.Server{
		Addr:         conf.Svr.Addr,
		ReadTimeout:  conf.Svr.ReadTimeout,
		WriteTimeout: conf.Svr.WriteTimeout,
		IdleTimeout:  conf.Svr.IdleTimeout,
		Handler:      mux,
	}
	return svr, nil
}
