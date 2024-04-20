package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/newrelic/go-agent/v3/newrelic"

	"go-playground/internal/transportlayer/rest/handler"
	"go-playground/internal/transportlayer/rest/middleware"
	"go-playground/pkg/nrmux"
	"go-playground/pkg/nrslog"
	"go-playground/pkg/problemdetails"
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
	mux.Handle("GET /reply/{name}", middleware.Recover(handler.ReplyHandler))

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
