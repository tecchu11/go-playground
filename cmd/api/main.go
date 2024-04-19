package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	svr, err := Initialize()
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
