package handler

import (
	"fmt"
	"go-playground/internal/transportlayer/rest/middleware"
	"go-playground/pkg/renderer"
	"log/slog"
	"net/http"
)

type HelloHandler interface {
	GetName() http.HandlerFunc
}

type helloHandler struct {
	rj renderer.JSON
}

func NewHelloHandler(failure renderer.JSON) HelloHandler {
	return &helloHandler{failure}
}

func (handler *helloHandler) GetName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := middleware.CurrentUser(r.Context())
		if err != nil {
			slog.ErrorContext(r.Context(), "Authenticated User does not exist in the request context", slog.String("path", r.URL.Path))
			handler.rj.Failure(w, r, http.StatusUnauthorized, "Request With No Authentication", "Request token was not found in your request header")
			return
		}
		slog.InfoContext(r.Context(), "hello !!", slog.String("user", user.Name))
		message := fmt.Sprintf("Hello %s!! You have %s role.", user.Name, user.Role.String())
		handler.rj.Success(w, 200, &HelloResponse{message})
	}
}

type HelloResponse struct {
	Message string `json:"message"`
}
