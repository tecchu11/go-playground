package service

import (
	"go-playground/configs"

	"github.com/newrelic/go-agent/v3/newrelic"

	"go-playground/internal/transportlayer/rest/handler"
	"go-playground/internal/transportlayer/rest/middleware"
	"go-playground/pkg/renderer"
	"go-playground/pkg/router"
)

// New returns configured router.
func New(prop *configs.ApplicationProperties, nrApp *newrelic.Application) *router.Router {
	// init misc
	jsonResponse := renderer.NewJSON(middleware.RequestID)
	preAuth := middleware.PreAuthenticatedUsers(make(map[string]middleware.AuthUser))
	for _, v := range prop.AuthConfigs {
		r, _ := middleware.UserRoleString(v.RoleStr) // ignore error
		preAuth[v.Key] = middleware.AuthUser{Name: v.Name, Role: r}
	}
	// init middleware
	newrelicMiddleware := middleware.NewrelicTxn(nrApp)
	recoverMiddleware := middleware.Recover(jsonResponse)
	generalAuth := preAuth.Auth(jsonResponse, map[middleware.UserRole]struct{}{middleware.Admin: {}, middleware.User: {}})
	// init handler
	helloHandler := handler.NewHelloHandler(jsonResponse)
	// init router
	r := router.New(
		router.Middlewares(newrelicMiddleware, recoverMiddleware),
		router.NotFoundHandlerPattern("NotFoundHandler"),
		router.MethodNotAllowedPattern("MethodNotAllowedHandler"),
	)
	// init routing
	r.Handle("GET /statuses", handler.StatusHandler(jsonResponse))
	r.Handle("GET /hello", generalAuth(helloHandler.GetName()))
	return r
}
