package handler_test

import (
	"go-playground/internal/transportlayer/rest/handler"
	"net/http"
	"testing"
)

var mux = http.NewServeMux()

func TestMain(m *testing.M) {
	mux.HandleFunc("GET /health", handler.HealthCheck)
	mux.HandleFunc("GET /reply/{name}", handler.ReplyHandler)
	m.Run()
}
