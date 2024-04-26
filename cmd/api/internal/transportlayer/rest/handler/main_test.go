package handler_test

import (
	"go-playground/cmd/api/internal/transportlayer/rest/handler"
	"net/http"
	"testing"
)

var mux = http.NewServeMux()

func TestMain(m *testing.M) {
	mux.HandleFunc("GET /health", handler.HealthCheck)
	mux.HandleFunc("GET /reply/{name}", handler.ReplyHandler)
	m.Run()
}
