package router_test

import (
	"context"
	"fmt"
	"go-playground/pkg/router"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithMiddleware(t *testing.T) {
	mid := router.Middleware(func(next http.Handler) http.Handler {
		return next // noop
	})
	optFn := router.WithMiddleware(mid)
	var router router.Router
	optFn(&router)
	require.NotNil(t, router.Middleware())
}

func TestWith404Func(t *testing.T) {
	fn := func(r *http.Request) ([]byte, error) {
		return nil, nil
	}
	optFn := router.With404Func(fn)
	var router router.Router
	optFn(&router)
	require.NotNil(t, router.Marshal404())
}

func TestWith405Func(t *testing.T) {
	fn := func(r *http.Request) ([]byte, error) {
		return nil, nil
	}
	optFn := router.With405Func(fn)
	var router router.Router
	optFn(&router)
	require.NotNil(t, router.Marshal405())
}

func TestNewWithoutOption(t *testing.T) {
	router := router.New()
	require.Equal(t, http.NewServeMux(), router.ServeMux)
	require.NotNil(t, router.Marshal404())
	require.NotNil(t, router.Marshal405())
	require.Nil(t, router.Middleware())
}

func TestNewWithOption(t *testing.T) {
	router := router.New(router.WithMiddleware(func(next http.Handler) http.Handler { return next }))
	require.NotNil(t, router.Middleware())
}

func TestRouter_ServeHTTP_WithDefault(t *testing.T) {
	tests := map[string]struct {
		w    *httptest.ResponseRecorder
		r    *http.Request
		code int
		body []byte
	}{
		"GET /index is 200": {
			w:    httptest.NewRecorder(),
			r:    httptest.NewRequest("GET", "https://example.com/index", nil),
			code: 200,
			body: []byte("this is /index"),
		},
		"POST /index is 405": {
			w:    httptest.NewRecorder(),
			r:    httptest.NewRequest("POST", "https://example.com/index", nil),
			code: 405,
			body: nil,
		},
		"GET /idx is 404": {
			w:    httptest.NewRecorder(),
			r:    httptest.NewRequest("POST", "https://example.com/idx", nil),
			code: 404,
			body: nil,
		},
		"RequestURI is asterisk": {
			w:    httptest.NewRecorder(),
			r:    httptest.NewRequest("GET", "*", nil),
			code: 400,
			body: nil,
		},
	}
	router := router.New()
	index := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		w.Write([]byte(fmt.Sprintf("this is %s", path)))
	})
	router.HandleFunc("GET /index", index)
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			t.Parallel()
			router.ServeHTTP(v.w, v.r)
			require.Equal(t, v.code, v.w.Code)
			require.Equal(t, v.body, v.w.Body.Bytes())
		})
	}
}

func TestRouter_ServeHTTP_WithCustom(t *testing.T) {
	tests := map[string]struct {
		w        *httptest.ResponseRecorder
		r        *http.Request
		code     int
		body     []byte
		reqCount int32
	}{
		"GET /index is 200": {
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest("GET", "https://example.com/index", nil),
			code:     200,
			body:     []byte("this is /index"),
			reqCount: 1,
		},
		"POST /index is 405": {
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest("POST", "https://example.com/index", nil),
			code:     405,
			body:     nil,
			reqCount: 1,
		},
		"GET /idx is 404": {
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest("POST", "https://example.com/idx", nil),
			code:     404,
			body:     nil,
			reqCount: 1,
		},
		"RequestURI is asterisk": {
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest("GET", "*", nil),
			code:     400,
			body:     nil,
			reqCount: 0,
		},
	}
	var reqCount int32
	router := router.New(
		router.WithMiddleware(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				atomic.AddInt32(&reqCount, 1)
				next.ServeHTTP(w, r)
			})
		}),
	)
	index := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		w.Write([]byte(fmt.Sprintf("this is %s", path)))
	})
	router.HandleFunc("GET /index", index)
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			router.ServeHTTP(v.w, v.r)
			require.Equal(t, v.code, v.w.Code)
			require.Equal(t, v.body, v.w.Body.Bytes())
			require.Equal(t, v.reqCount, reqCount)
			atomic.StoreInt32(&reqCount, 0) // reset
		})
	}
}

func TestPattern(t *testing.T) {
	tests := map[string]struct {
		ctx      context.Context
		expected any
	}{
		"with routing pattern": {
			ctx:      context.WithValue(context.Background(), router.PatternContextKey{}, "GET /index"),
			expected: "GET /index",
		},
		"with routing missing pattern": {
			ctx:      context.Background(),
			expected: "MissingRoutingPattern",
		},
		"with routing missing pattern string": {
			ctx:      context.WithValue(context.Background(), router.PatternContextKey{}, []byte("GET /index")),
			expected: "MissingRoutingPatternString",
		},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			req, err := http.NewRequestWithContext(v.ctx, "GET", "https://example.com/index", nil)
			require.NoError(t, err)
			actual := router.Pattern(req)
			require.Equal(t, v.expected, actual)
		})
	}
}
