package main

import (
	"fmt"
	"go-playground/config"
	"go-playground/pkg/presentation"
	"go-playground/pkg/presentation/auth"
	"go-playground/pkg/presentation/handler"
	"go-playground/pkg/presentation/middleware"
	"log"
	"net/http"
	"os"

	"go.uber.org/zap"
)

var (
	configFile string
)

func init() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		log.Fatal("APP_ENV is not stored in evirioment variables")
	}
	configFile = fmt.Sprintf("config-%s.json", env)
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("zap logger is failed to init because of %s", err)
	}

	prop := config.NewPropertiesLoader(logger).Load(configFile)
	appLogger := logger.With(zap.String("appName", prop.AppName))
	appLogger.Info("Success to load properties")

	// initialize middleware
	authMid := middleware.NewAuthenticationMiddleWare(appLogger, auth.NewAutheticatonManager(prop.AuthConfigs))
	authenticatedCompostionMiddleware := middleware.Composite(authMid.Handle)
	// initialize handler
	health := handler.NewHealthHandler(appLogger).GetStatus()
	hello := handler.NewHelloHandler(appLogger).GetName()

	mux := presentation.NewMuxBuilder().
		SetHadler("/health", health).
		SetHadler("/hello", authenticatedCompostionMiddleware(hello)).
		Build()

	appLogger.Info("Server started ---(ﾟ∀ﾟ)---!!!")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		appLogger.Fatal("Failed to start server", zap.Error(err))
	}
}
