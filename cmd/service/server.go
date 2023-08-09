package service

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go-playground/configs"
	"go-playground/internal/transport_layer/rest/handler"
	"go-playground/internal/transport_layer/rest/middleware"
	"go-playground/internal/transport_layer/rest/preauth"
	"go-playground/pkg/renderer"
	"go.uber.org/zap"
)

var (
	ErrInvalidEnv = errors.New("invalid env")
	ErrInitNRApp  = errors.New("failed to init newrelic app")
)

// New returns configured chi.mux.
func New(env string, logger *zap.Logger, prop *configs.ApplicationProperties) (*chi.Mux, error) {
	if env == "" {
		return nil, ErrInvalidEnv
	}
	failure := renderer.NewFailure(middleware.RequestID)
	mux := chi.NewRouter()
	mux.MethodNotAllowed(handler.MethodNotAllowedHandler(failure))
	mux.NotFound(handler.NotFoundHandler(failure))
	nrApp, err := newrelicApp(env)
	if err != nil {
		return nil, errors.Join(ErrInitNRApp, err)
	}
	mux.Use(middleware.NewrelicTxn(nrApp))
	mux.Use(middleware.Recover(logger, failure))

	authMid := middleware.Authenticator(logger, preauth.NewAutheticatonManager(prop.AuthConfigs), failure)
	hello := handler.NewHelloHandler(logger, failure).GetName()
	mux.Route("/statuses", func(r chi.Router) {
		r.Get("/", handler.StatusHandler())
	})
	mux.Route("/hello", func(r chi.Router) {
		r.Use(authMid)
		r.Get("/", hello)
	})
	return mux, nil
}

func newrelicApp(env string) (*newrelic.Application, error) {
	if env == "local" {
		return nil, nil
	}
	return newrelic.NewApplication(newrelic.ConfigFromEnvironment())
}
