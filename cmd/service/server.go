package service

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go-playground/configs"
	handler2 "go-playground/internal/interactor/rest/handler"
	middleware2 "go-playground/internal/interactor/rest/middleware"
	"go-playground/internal/interactor/rest/preauth"
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
	mux := chi.NewRouter()
	mux.MethodNotAllowed(handler2.MethodNotAllowedHandler())
	mux.NotFound(handler2.NotFoundHandler())
	nrApp, err := newrelicApp(env)
	if err != nil {
		return nil, errors.Join(ErrInitNRApp, err)
	}
	mux.Use(middleware2.NewrelicTxn(nrApp))
	mux.Use(middleware2.Recover(logger))

	authMid := middleware2.Autheticator(logger, preauth.NewAutheticatonManager(prop.AuthConfigs))
	hello := handler2.NewHelloHandler(logger).GetName()
	mux.Route("/statuses", func(r chi.Router) {
		r.Get("/", handler2.StatusHandler())
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
