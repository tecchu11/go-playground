package middleware_test

import (
	"go-playground/internal/transportlayer/rest/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/require"
)

func TestNewRelicTxn(t *testing.T) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("test-local"),
		newrelic.ConfigLicense("0000000000000000000000000000000000000000"),
	)
	require.NoError(t, err)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txn := newrelic.FromContext(r.Context())
		require.NotNil(t, txn)
		w.Write([]byte("ok"))
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "https://example.com/ok", nil)
	middleware.NewrelicTxn(app)(handler).ServeHTTP(w, r)
	require.Equal(t, 200, w.Code)
	require.Equal(t, []byte("ok"), w.Body.Bytes())
}

func TestNewRelicTxn_AppIsNil(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txn := newrelic.FromContext(r.Context())
		require.Nil(t, txn)
		w.Write([]byte("ok"))
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "https://example.com/ok", nil)
	middleware.NewrelicTxn(nil)(handler).ServeHTTP(w, r)
	require.Equal(t, 200, w.Code)
	require.Equal(t, []byte("ok"), w.Body.Bytes())
}
