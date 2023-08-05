package main

import (
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

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		log.Fatal("Missing APP_ENV in enviroment variables")
	}
	prop, err := config.Load(env)
	if err != nil {
		log.Fatalf("Failed to load configuration with %s because %v", env, err)
	}
	logger, err := zapLogger(env, prop.AppName)
	if err != nil {
		log.Fatal("Failed to init zap logger", err)
	}
	app, err := newrelicApp(env)
	if err != nil {
		logger.Fatal("failed to init newrelic application", zap.Error(err))
	}

	// initialize middleware
	recoverMid := middleware.NewRecoverMiddleWare(logger)
	authMid := middleware.NewAuthenticationMiddleWare(logger, preauth.NewAutheticatonManager(prop.AuthConfigs))
	newrelicTxnMid := middleware.NewNewrelicTransactionMidleware(app)
	// initialize handler
	hello := handler.NewHelloHandler(logger).GetName()

	mux := chi.NewRouter()
	mux.MethodNotAllowed(handler.MethodNotAllowedHandler())
	mux.NotFound(handler.NotFoundHandler())
	mux.Use(newrelicTxnMid.Handle)
	mux.Use(recoverMid.Handle)
	mux.Route("/statuses", func(r chi.Router) {
		r.Get("/", handler.StatusHandler())
	})
	mux.Route("/hello", func(r chi.Router) {
		r.Use(authMid.Handle)
		r.Get("/", hello)
	})

	logger.Info("Server started ---(ﾟ∀ﾟ)---!!!")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
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

// newrelicApp init newrelic.Application. When env is local, app is returned nil.
func newrelicApp(env string) (*newrelic.Application, error) {
	if env == "local" {
		return nil, nil
	}
	return newrelic.NewApplication(
		newrelic.ConfigFromEnvironment(),
	)
}
