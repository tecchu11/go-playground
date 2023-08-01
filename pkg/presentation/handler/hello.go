package handler

import (
	"fmt"
	"go-playground/pkg/lib/render"
	"go-playground/pkg/presentation/middleware"
	"go-playground/pkg/presentation/model"
	"net/http"

	"go.uber.org/zap"
)

type HelloHandler interface {
	GetName() http.HandlerFunc
}

type helloHandler struct {
	logger *zap.Logger
}

func NewHelloHandler(logger *zap.Logger) HelloHandler {
	return &helloHandler{logger}
}

func (handler *helloHandler) GetName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := middleware.GetAutenticatedUser(r.Context())
		if err != nil {
			handler.logger.Error("Authenticated User does not exsist in the request context", zap.String("path", r.URL.Path))
			render.Unauthorized(w, "No token was found for your request", r.URL.Path)
			return
		}
		message := fmt.Sprintf("Hello %s!! You have %s role.", user.Name, user.Role.String())
		render.Ok(w, &model.HelloResponse{Message: message})
	}

}
