package main

import (
	"go-playground/pkg/nrmux"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var svr *http.Server

func TestMain(m *testing.M) {
	os.Setenv("APP_ENV", "local")
	os.Setenv("NEW_RELIC_LICENSE_KEY", "0000000000000000000000000000000000000000")
	os.Setenv("NEW_RELIC_APP_NAME", "test-local")
	var err error
	svr, err = Initialize()
	if err != nil {
		panic(err)
	}
	m.Run()
}

func TestInitialize(t *testing.T) {
	require.Equal(t, ":8080", svr.Addr)
	require.Equal(t, 10*time.Second, svr.ReadTimeout)
	require.Equal(t, 10*time.Second, svr.WriteTimeout)
	require.Equal(t, 120*time.Second, svr.IdleTimeout)
	require.NotNil(t, svr.Handler)
}

func TestRouting(t *testing.T) {
	tests := map[string]struct {
		r        *http.Request
		expected string
	}{
		"GET /health":       {r: httptest.NewRequest("GET", "https://example.com/health", nil), expected: "GET /health"},
		"GET /reply/{name}": {r: httptest.NewRequest("GET", "https://example.com/reply/Cummerata", nil), expected: "GET /reply/{name}"},
	}
	mux, ok := svr.Handler.(*nrmux.NRServeMux)
	require.True(t, ok)

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			_, actual := mux.Handler(v.r)

			require.Equal(t, v.expected, actual)
		})
	}

}
