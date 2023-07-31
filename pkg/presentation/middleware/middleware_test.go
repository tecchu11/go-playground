package middleware_test

import (
	"fmt"
	"go-playground/pkg/presentation/middleware"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type caller struct {
	calls []string
}

func (c *caller) log(str string) {
	c.calls = append(c.calls, str)
}

func TestComposite(t *testing.T) {
	c := caller{}
	first := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.log("first")
			handler.ServeHTTP(w, r)
		})
	}
	second := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.log("second")
			handler.ServeHTTP(w, r)
		})
	}
	third := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c.log("third")
			handler.ServeHTTP(w, r)
		})
	}
	fx := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "last")
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://example.com", nil)

	middleware.Composite(first, second, third)(fx).ServeHTTP(w, r)

	if !reflect.DeepEqual(c.calls, []string{"first", "second", "third"}) || w.Body.String() != "last" {
		t.Errorf("unexpected caller result. caller is (%v) and response is (%v)", c.calls, w.Body.String())
	}
}