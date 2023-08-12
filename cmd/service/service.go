package service

import (
	"github.com/go-chi/chi/v5"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go-playground/configs"

	"go-playground/internal/transport_layer/rest/handler"
	"go-playground/internal/transport_layer/rest/middleware"
	"go-playground/pkg/renderer"
	"go.uber.org/zap"
)

// New returns configured chi.mux.
func New(env string, logger *zap.Logger, prop *configs.ApplicationProperties, nrApp *newrelic.Application) *chi.Mux {
	// init misc
	failure := renderer.NewFailure(middleware.RequestID)
	preAuth := middleware.PreAuthenticatedUsers(make(map[string]middleware.AuthUser))
	for _, v := range prop.AuthConfigs {
		r, _ := middleware.UserRoleString(v.RoleStr) // ignore error
		preAuth[v.Key] = middleware.AuthUser{Name: v.Name, Role: r}
	}
	// init middleware
	generalAuth := preAuth.Auth(logger, failure, map[middleware.UserRole]struct{}{middleware.Admin: {}, middleware.User: {}})
	// init handler
	helloHandler := handler.NewHelloHandler(logger, failure)
	// init mux
	mux := chi.NewRouter()
	mux.MethodNotAllowed(handler.MethodNotAllowedHandler(failure))
	mux.NotFound(handler.NotFoundHandler(failure))
	mux.Use(middleware.NewrelicTxn(nrApp))
	mux.Use(middleware.Recover(logger, failure))
	// init routing
	mux.Route("/statuses", func(r chi.Router) {
		r.Get("/", handler.StatusHandler())
	})
	mux.Route("/hello", func(r chi.Router) {
		r.Use(generalAuth)
		r.Get("/", helloHandler.GetName())
	})
	return mux
}
