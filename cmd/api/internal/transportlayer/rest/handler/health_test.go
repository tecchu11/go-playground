package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "https://example.com/health", nil)

	mux.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"status":"ok"}`, w.Body.String())
}
