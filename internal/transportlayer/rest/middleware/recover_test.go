package middleware_test

import (
	"go-playground/internal/transportlayer/rest/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Recover(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("always panic")
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "https://example.com/panic", nil)
	middleware.Recover(handler).ServeHTTP(w, r)
	require.Equal(t, 500, w.Code)
	require.Equal(t, "application/problem+json", w.Header().Get("Content-Type"))
	require.JSONEq(t, `{"type":"about:blank","title":"Internal Server Error","status":500,"detail":"Unexpected error was happened. Please report this error you have checked.","instance":"/panic"}`, w.Body.String())
}

func Test_Recover_PanicWithErrAbortHandler(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic(http.ErrAbortHandler)
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "https://example/com/abort", nil)
	defer func() {
		err := recover()
		require.Equal(t, http.ErrAbortHandler, err)
	}()
	middleware.Recover(handler).ServeHTTP(w, r)
}
