package nrmux_test

import (
	"fmt"
	"go-playground/pkg/nrmux"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/stretchr/testify/require"
)

func TestWithMarshalJSON404(t *testing.T) {
	var opt nrmux.Option
	nrmux.WithMarshalJSON404(func(r *http.Request) ([]byte, error) {
		return nil, nil
	})(&opt)
	require.NotNil(t, opt.MarshalJSON404())
}

func TestWithMarshalJSON405(t *testing.T) {
	var opt nrmux.Option
	nrmux.WithMarshalJSON405(func(r *http.Request) ([]byte, error) {
		return nil, nil
	})(&opt)
	require.NotNil(t, opt.MarshalJSON405())
}

func TestNew(t *testing.T) {
	tests := map[string]struct {
		opts             []nrmux.OptionFunc
		expectedApp      *newrelic.Application
		expectedServeMux *http.ServeMux
	}{
		"default": {
			expectedApp:      app,
			expectedServeMux: http.NewServeMux(),
		},
		"with args": {
			opts: []nrmux.OptionFunc{
				nrmux.WithMarshalJSON404(func(r *http.Request) ([]byte, error) {
					return nil, nil
				}),
				nrmux.WithMarshalJSON405(func(r *http.Request) ([]byte, error) {
					return nil, nil
				}),
			},
			expectedApp:      app,
			expectedServeMux: http.NewServeMux(),
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			mux := nrmux.New(app, v.opts...)

			require.Equal(t, mux.ServeMux, v.expectedServeMux)
			require.Equal(t, mux.App(), v.expectedApp)
			require.NotNil(t, mux.Wrap())
			require.NotNil(t, mux.Unwrap())
		})
	}
}

func TestNRServeMux_Handle(t *testing.T) {
	mux := nrmux.New(app)
	ok := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"msg":"ok"}`))
	})
	mux.Handle("GET /ok", ok)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "https://example.com/ok", nil)
	h, p := mux.ServeMux.Handler(httptest.NewRequest("GET", "https://example.com/ok", nil))
	h.ServeHTTP(w, req)

	require.Equal(t, "GET /ok", p)
	require.Equal(t, 200, w.Code)
	require.Equal(t, []byte(`{"msg":"ok"}`), w.Body.Bytes())
}

func TestNRServeMux_HandleFunc(t *testing.T) {
	mux := nrmux.New(app)
	ok := func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"msg":"ok"}`))
	}
	mux.HandleFunc("GET /ok", ok)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "https://example.com/ok", nil)
	h, p := mux.ServeMux.Handler(httptest.NewRequest("GET", "https://example.com/ok", nil))
	h.ServeHTTP(w, req)

	require.Equal(t, "GET /ok", p)
	require.Equal(t, 200, w.Code)
	require.Equal(t, []byte(`{"msg":"ok"}`), w.Body.Bytes())
}

func TestNRServeMux_ServeHTTP(t *testing.T) {
	mux := nrmux.New(app)
	ok := func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"msg":"ok"}`))
	}
	hi := func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		w.Write([]byte(fmt.Sprintf("{\"msg\":\"Hi %s\"}", name)))
	}
	mux.HandleFunc("GET /ok", ok)
	mux.HandleFunc("GET /hi/{name}", hi)

	tests := map[string]struct {
		w *httptest.ResponseRecorder
		r *http.Request

		expectedCode int
		expectedBody []byte
	}{
		"GET request to /ok": {
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest("GET", "https://example.com/ok", nil),
			expectedCode: 200,
			expectedBody: []byte(`{"msg":"ok"}`),
		},
		"GET request to /404": {
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest("GET", "https://example.com/404", nil),
			expectedCode: 404,
			expectedBody: nil,
		},
		"POST request to /ok": {
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest("POST", "https://example.com/ok", nil),
			expectedCode: 405,
			expectedBody: nil,
		},
		"GET request to /hi/{name}": {
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest("GET", "https://example.com/hi/Auer", nil),
			expectedCode: 200,
			expectedBody: []byte(`{"msg":"Hi Auer"}`),
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			mux.ServeHTTP(v.w, v.r)

			require.Equal(t, v.expectedCode, v.w.Code)
			require.Equal(t, v.expectedBody, v.w.Body.Bytes())
		})
	}
}
