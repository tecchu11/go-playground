package main

import (
	"context"
	"errors"
	"go-playground/cmd/service"
	"go-playground/config"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		log.Fatal("Missing APP_ENV in environment variables")
	}
	prop, err := config.Load(env)
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
	mux, err := service.New(env, logger, prop)
	if err != nil {
		log.Fatal("Failed to init service mux", err)
	}

	svr := &http.Server{Addr: ":8080", Handler: mux}
	logger.Info("Server starting ---(ﾟ∀ﾟ)---!!!")
	go func() {
		if err := svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("failed to start up server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Termination signal was caught. So attempt to terminate server with gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func(cancel context.CancelFunc) {
		logger.Info("Bye!!")
		cancel()
	}(cancel)
	if err := svr.Shutdown(ctx); err != nil {
		logger.Fatal("Failed to terminate server with gracefully. So force terminating ...", zap.Error(err))
	}
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
