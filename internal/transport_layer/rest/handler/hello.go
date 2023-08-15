package handler

import (
	"fmt"
	"go-playground/internal/transport_layer/rest/middleware"
	"go-playground/pkg/renderer"
	"net/http"

	"go.uber.org/zap"
)

type HelloHandler interface {
	GetName() http.HandlerFunc
}

type helloHandler struct {
	logger *zap.Logger
	rj     renderer.JSON
}

func NewHelloHandler(logger *zap.Logger, failure renderer.JSON) HelloHandler {
	return &helloHandler{logger, failure}
}

func (handler *helloHandler) GetName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := middleware.CurrentUser(r.Context())
		if err != nil {
			handler.logger.Error("Authenticated User does not exist in the request context", zap.String("path", r.URL.Path))
			handler.rj.Failure(w, r, http.StatusUnauthorized, "Request With No Authentication", "Request token was not found in your request header")
			return
		}
		message := fmt.Sprintf("Hello %s!! You have %s role.", user.Name, user.Role.String())
		handler.rj.Success(w, 200, &HelloResponse{message})
	}
}

type HelloResponse struct {
	Message string `json:"message"`
}
