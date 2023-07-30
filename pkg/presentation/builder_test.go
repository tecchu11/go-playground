package presentation_test

import (
	"fmt"
	"go-playground/pkg/presentation"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMuxBuidler_SetHadler(t *testing.T) {
	builder := presentation.NewMuxBuilder()
	builder.SetHadler("/foos", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/foos" {
			t.Errorf("unexpected path(%v) was routed", r.URL.Path)
		}
		_, _ = fmt.Fprint(w, "test")
	}))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://example.com/foos", nil)
	builder.Build().ServeHTTP(w, r)
	if string(w.Body.Bytes()) != "test" {
		t.Errorf("unexpected body(%v) was recived", string(w.Body.Bytes()))
	}
}

func TestMuxBuidler_SetHandlerFunc(t *testing.T) {
	builder := presentation.NewMuxBuilder()
	builder.SetHandlerFunc("/bars", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/bars" {
			t.Errorf("unexpected path(%v) was routed", r.URL.Path)
		}
		_, _ = fmt.Fprint(w, "test")
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "http://example.com/bars", nil)
	builder.Build().ServeHTTP(w, r)
	if string(w.Body.Bytes()) != "test" {
		t.Errorf("unexpected body(%v) was recived", string(w.Body.Bytes()))
	}
}
