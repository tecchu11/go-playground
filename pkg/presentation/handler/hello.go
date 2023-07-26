package handler

import (
	"fmt"
	"go-playground/pkg/presentation/auth"
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

func NewHelloHandler(log *zap.Logger) HelloHandler {
	return &helloHandler{logger: log}
}

func (handler *helloHandler) GetName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := auth.GetAuthUser(r.Context())
		if err != nil {
			handler.logger.Error("Authenticated User does not exsist in the request context", zap.String("path", r.URL.Path))
			Unauthorized(w, "No token was found for your request", r.URL.Path)
			return
		}
		message := fmt.Sprintf("Hello %s!! You have %s role.", user.Name, user.Role.String())
		Ok(w, &model.HelloResponse{Message: message})
	}

}
