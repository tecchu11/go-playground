package handler

import (
	"encoding/json"
	"go-playground/pkg/presentation/model"
	"net/http"

	"go.uber.org/zap"
)

type HelloHandler interface {
	GetName() http.HandlerFunc
}

type helloHandler struct {
	log *zap.Logger
}

func NewHelloHandler(log *zap.Logger) HelloHandler {
	return &helloHandler{log: log}
}

func (hh *helloHandler) GetName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uname := r.Context().Value("user")
		if uname == nil {
			hh.log.Error("User name does not exsist in context", zap.String("path", r.URL.Path))
			p := model.NewProblemDetail("You had requested invalid token", r.URL.Path, http.StatusUnauthorized)
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(p); err != nil {
				hh.log.Error("Failed to write response", zap.Error(err), zap.String("path", r.URL.Path))
			}
			return
		}
		hello := model.HelloResponse{Message: "Hello !!", Name: uname.(string)}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(hello); err != nil {
			hh.log.Error("Failed to write response", zap.Error(err), zap.String("path", r.URL.Path))
		}
	}

}
