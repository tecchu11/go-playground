package main

import (
	"fmt"
	"go-playground/config"
	"go-playground/pkg/presentation/handler"
	"go-playground/pkg/presentation/middleware"
	"go-playground/pkg/presentation/preauth"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
)

var (
	env  string
	prop *config.Properties
)

func init() {
	env = os.Getenv("APP_ENV")
	if env == "" {
		log.Fatal("APP_ENV is not stored in evirioment variables")
	}
	var err error
	prop, err = config.LoadConfigWith(fmt.Sprintf("config-%s.json", env))
	if err != nil {
		log.Fatal("failed to load configuration", err)
	}
}

func main() {
	logger, _ := zap.NewProduction()
	appLogger := logger.
		With(zap.String("appName", prop.AppName)).
		With(zap.String("env", env))

	app, err := newrelicApp()
	if err != nil {
		appLogger.Fatal("failed to init newrelic application", zap.Error(err))
	}

	// initialize middleware
	authMid := middleware.NewAuthenticationMiddleWare(appLogger, preauth.NewAutheticatonManager(prop.AuthConfigs))
	newrelicTxnMid := middleware.NewNewrelicTransactionMidleware(app)
	// initialize handler
	hello := handler.NewHelloHandler(appLogger).GetName()

	mux := chi.NewRouter()
	mux.MethodNotAllowed(handler.MethodNotAllowedHandler())
	mux.NotFound(handler.NotFoundHandler())
	mux.Use(newrelicTxnMid.Handle)
	mux.Route("/statuses", func(r chi.Router) {
		r.Get("/", handler.StatusHandler())
	})
	mux.Route("/hello", func(r chi.Router) {
		r.Use(authMid.Handle)
		r.Get("/", hello)
	})

	appLogger.Info("Server started ---(ﾟ∀ﾟ)---!!!")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		appLogger.Fatal("Failed to start server", zap.Error(err))
	}
}

func newrelicApp() (*newrelic.Application, error) {
	if env == "local" {
		return nil, nil
	}
	return newrelic.NewApplication(
		newrelic.ConfigFromEnvironment(),
	)
}
