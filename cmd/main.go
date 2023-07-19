package main

import (
	"go-playground/config"
	"go-playground/pkg/presentation"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	log := config.NewLogger()
	mux := presentation.BuildMux(http.NewServeMux(), log)
	log.Info("Server started")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
