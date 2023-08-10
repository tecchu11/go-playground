package main

import (
	"context"
	"errors"
	"go-playground/cmd/service"
	"go-playground/configs"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		log.Fatal("Missing APP_ENV in environment variables")
	}
	prop, err := configs.Load(env)
	if err != nil {
		log.Fatalf("Failed to load configuration with %s because %v", env, err)
	}
	logger, err := zapLogger(env, prop.AppName)
	if err != nil {
		log.Fatal("Failed to init zap logger", err)
	}
	defer func(logger *zap.Logger) {
		_ = logger.Sync() // ignore sync error.
	}(logger)
	nrApp, err := newrelicApp(env)
	if err != nil {
		log.Fatal("Failed to init newrelic Application", err)
	}

	mux := service.New(env, logger, prop, nrApp)
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

		logger.Info("We received an interrupt signal,so attempt to shut down with gracefully")
		ctx, cancel := context.WithTimeout(context.Background(), prop.ServerConfig.GraceTimeout)
		defer func(logger *zap.Logger) {
			logger.Info("Bye!!")
			cancel()
		}(logger)
		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("Failed to terminate server with gracefully. So force terminating ...", zap.Error(err))
		}
		close(idleConnsClosed)
	}()

	logger.Info("Server starting ---(ﾟ∀ﾟ)---!!!")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("Failed to start up server", zap.Error(err))
	}
	<-idleConnsClosed
}

// zapLogger init development or production logger with zap filed of env, appName.
func zapLogger(env string, appName string) (*zap.Logger, error) {
	opt := zap.Fields(
		zap.String("appName", appName),
		zap.String("env", env),
	)
	if env == "local" {
		return zap.NewDevelopment(opt)
	}
	return zap.NewProduction(opt)
}

// newrelicApp init *newrelic.Application. If env is local, *newrelic.Application is nil.
func newrelicApp(env string) (*newrelic.Application, error) {
	if env == "local" {
		return nil, nil
	}
	return newrelic.NewApplication(newrelic.ConfigFromEnvironment())
}
