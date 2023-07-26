package main

import (
	"fmt"
	"go-playground/config"
	"go-playground/pkg/presentation"
	"go-playground/pkg/presentation/auth"
	"go-playground/pkg/presentation/handler"
	"go-playground/pkg/presentation/middleware"
	"net/http"
	"os"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}

	configLocation := fmt.Sprintf("../config/config-%s.json", "local")
	prop := config.NewPropertiesLoader(logger).Load(configLocation)
	appLogger := logger.With(zap.String("appName", prop.AppName))
	appLogger.Info("Success to load properties")

	headMid := middleware.NewHeadMiddleWare()
	ctxMid := middleware.NewContextMiddleWare()
	authMid := middleware.NewAuthMiddleWare(appLogger, &auth.AuthenticationManager{Configs: prop.AuthConfigs})
	health := handler.NewHealthHandler(appLogger).GetStatus()
	hello := handler.NewHelloHandler(appLogger).GetName()

	mux := presentation.NewMuxBuilder().
		SetHadler("/health", middleware.Composite(headMid.Handle)(health)).
		SetHadler("/hello", middleware.Composite(headMid.Handle, ctxMid.Handle, authMid.Handle)(hello)).
		Build()

	appLogger.Info("Server started ---(ﾟ∀ﾟ)---!!!")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		appLogger.Fatal("Failed to start server", zap.Error(err))
	}
}
