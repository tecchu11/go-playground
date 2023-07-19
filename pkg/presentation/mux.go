package presentation

import (
	"go-playground/pkg/presentation/handler"
	"go-playground/pkg/presentation/middleware"
	"net/http"

	"go.uber.org/zap"
)

func BuildMux(mux *http.ServeMux, log *zap.Logger) *http.ServeMux {
	mux.HandleFunc("/health", handler.NewHealthHandler(log).GetStatus())
	mux.HandleFunc(
		"/hello",
		middleware.ContextMiddleWareFunc(
			middleware.AuthMiddleWareFunc(
				handler.NewHelloHandler(log).GetName(),
				log,
			),
		),
	)
	return mux
}
