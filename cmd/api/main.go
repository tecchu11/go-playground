package main

import (
	"context"
	"errors"
	"fmt"
	"go-playground/cmd/api/service"
	"go-playground/configs"
	"go-playground/pkg/nrslog"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
	env := os.Getenv("APP_ENV")
	prop, err := configs.Load(env)
	if err != nil {
		panic(fmt.Sprintf("failed to load config with env %s: %v", env, err))
	}
	app, err := newrelic.NewApplication(newrelic.ConfigFromEnvironment())
	if err != nil {
		panic(fmt.Sprintf("failed to init newrelic app: %v", err))
	}
	nrHandler, err := nrslog.NewJSONHandler(app, &slog.HandlerOptions{AddSource: true})
	if err != nil {
		panic(err)
	}
	slog.SetDefault(slog.New(nrHandler))

	mux := service.New(prop, app)
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
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		slog.Info("We received an interrupt signal,so attempt to shutdown with gracefully")
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
