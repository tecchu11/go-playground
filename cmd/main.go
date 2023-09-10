package main

import (
	"context"
	"errors"
	"go-playground/cmd/service"
	"go-playground/configs"
	"go-playground/pkg/nrslog"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func init() {
	slog.SetDefault(
		slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})),
	)
}

func main() {
	env := os.Getenv("APP_ENV")
	prop, err := configs.Load(env)
	if err != nil {
		slog.Error("Failed to load configuration", slog.String("env", env), slog.String("error", err.Error()))
		panic(err)
	}
	nrApp, err := newrelicApp(env)
	if err != nil {
		slog.Error("Failed to init newrelic Application", slog.String("error", err.Error()))
		panic(err)
	}
	slog.SetDefault(slog.New(nrslog.New(nrApp)))

	mux := service.New(prop, nrApp)
	srv := &http.Server{
		Addr:         prop.ServerConfig.Address,
		ReadTimeout:  prop.ServerConfig.ReadTimeout,
		WriteTimeout: prop.ServerConfig.WriteTimeout,
		IdleTimeout:  prop.ServerConfig.IdleTimeout,
		Handler:      mux,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		slog.Info("We received an interrupt signal,so attempt to shut down with gracefully")
		ctx, cancel := context.WithTimeout(context.Background(), prop.ServerConfig.GraceTimeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("Failed to terminate server with gracefully. So force terminating ...", slog.String("error", err.Error()))
		}
		close(idleConnsClosed)
	}()

	slog.Info("Server starting ---(ﾟ∀ﾟ)---!!!")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Failed to start up server", slog.String("error", err.Error()))
		panic(err)
	}
	<-idleConnsClosed
	slog.Info("Bye!!")
}

// newrelicApp init *newrelic.Application. If env is local, *newrelic.Application is nil.
func newrelicApp(env string) (*newrelic.Application, error) {
	if env == "local" {
		return nil, nil
	}
	return newrelic.NewApplication(newrelic.ConfigFromEnvironment())
}
