package handler_test

import (
	"go-playground/cmd/api/internal/transportlayer/rest/handler"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "https://example.com/health", nil)

	handler.HealthCheck.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}
