package middleware_test

import (
	"go-playground/pkg/presentation/middleware"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContextMiddleWare_Handle(t *testing.T) {
	tests := map[string]struct {
		inputWriter  *httptest.ResponseRecorder
		inputRequest *http.Request
		expected     string
	}{
		"stored value in context is expected": {httptest.NewRecorder(), httptest.NewRequest("GET", "http://example.com/foos", nil), "/foos"},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				actual := r.Context().Value("path").(string)
				if actual != v.expected {
					t.Errorf("unexpected value(%s) was stored in context", actual)
				}
			})
			middleware.NewContextMiddleWare().Handle(next).ServeHTTP(v.inputWriter, v.inputRequest)
		})
	}
}
