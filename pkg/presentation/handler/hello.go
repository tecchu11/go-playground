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
			problem := model.NewProblemDetail("You had requested invalid token", r.URL.Path, http.StatusUnauthorized)
			JsonResponse(w, http.StatusUnauthorized, problem)
			return
		}
		message := fmt.Sprint("Hello!! ", user.Name, " role is ", user.Role.String())
		hello := model.HelloResponse{Message: message}
		JsonResponse(w, http.StatusOK, hello)
	}

}
