package middleware_test

import (
	"go-playground/pkg/presentation/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeadMiddleWare_Handle(t *testing.T) {
	tests := map[string]struct {
		inputWriter  *httptest.ResponseRecorder
		inputRequest *http.Request
		expected     string
	}{
		"Content-Type is appplication/json": {httptest.NewRecorder(), httptest.NewRequest("GET", "http://example.com/foos", nil), "application/json"},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				actual := w.Header().Get("Content-Type")
				if actual != v.expected {
					t.Errorf("Content-Type is unexpected value(%s)", actual)
				}
			})
			middleware.NewHeadMiddleWare().Handle(next).ServeHTTP(v.inputWriter, v.inputRequest)
		})
	}

}
