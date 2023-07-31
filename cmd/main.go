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

var env string

func init() {
	env = os.Getenv("APP_ENV")
	if env == "" {
		log.Fatal("APP_ENV is not stored in evirioment variables")
	}
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("zap logger is failed to init because of %s", err)
	}

	configLocation := fmt.Sprintf("../config/config-%s.json", env)
	prop := config.NewPropertiesLoader(logger).Load(configLocation)
	appLogger := logger.With(zap.String("appName", prop.AppName))
	appLogger.Info("Success to load properties")

	// initialze misc
	authCtxManager := middleware.AuthCtxManager{}

	// initialize middleware
	headMid := middleware.NewHeadMiddleWare()
	authMid := middleware.NewAuthMiddleWare(appLogger, auth.NewAutheticatonManager(prop.AuthConfigs))
	noAuthenticatedCompostionMiddleware := middleware.Composite(headMid.Handle)
	authenticatedCompostionMiddleware := middleware.Composite(headMid.Handle, authMid.Handle)
	// initialize handler
	health := handler.NewHealthHandler(appLogger).GetStatus()
	hello := handler.NewHelloHandler(appLogger, &authCtxManager).GetName()

	mux := presentation.NewMuxBuilder().
		SetHadler("/health", noAuthenticatedCompostionMiddleware(health)).
		SetHadler("/hello", authenticatedCompostionMiddleware(hello)).
		Build()

	appLogger.Info("Server started ---(ﾟ∀ﾟ)---!!!")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		appLogger.Fatal("Failed to start server", zap.Error(err))
	}
}
