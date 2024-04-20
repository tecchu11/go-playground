package nrmux_test

import (
	"errors"
	"go-playground/pkg/nrmux"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrapResponseWriter_Unwrap(t *testing.T) {
	in := httptest.NewRecorder()
	wrw := nrmux.WrapResponseWriter{ResponseWriter: in}
	require.Equal(t, in, wrw.Unwrap())
}

func TestWrapResponseWriter_WriteHeader(t *testing.T) {
	wrw := nrmux.WrapResponseWriter{ResponseWriter: httptest.NewRecorder()}
	wrw.WriteHeader(200)
	require.Equal(t, 200, wrw.Code())
}

func TestWrapResponseWriter_Write(t *testing.T) {
	marshal404 := func(*http.Request) ([]byte, error) {
		return []byte(`{"code":404}`), nil
	}
	marshal405 := func(*http.Request) ([]byte, error) {
		return []byte(`{"code":405}`), nil
	}
	marshalErr := func(*http.Request) ([]byte, error) {
		return nil, errors.New("marshal error")
	}
	tests := map[string]struct {
		w            *httptest.ResponseRecorder
		r            *http.Request
		code         int
		marshal404   func(*http.Request) ([]byte, error)
		marshal405   func(*http.Request) ([]byte, error)
		buf          []byte
		expectedCode int
		expectedBody []byte
	}{
		"status 404": {
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest("GET", "https://example.com/404", nil),
			code:         404,
			marshal404:   marshal404,
			marshal405:   marshal405,
			buf:          []byte("404 page not found"),
			expectedCode: 404,
			expectedBody: []byte(`{"code":404}`),
		},
		"status 404 with marshal error": {
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest("GET", "https://example.com/404", nil),
			code:         404,
			marshal404:   marshalErr,
			marshal405:   marshal405,
			buf:          []byte("404 page not found"),
			expectedCode: 404,
			expectedBody: nil,
		},
		"status 405": {
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest("GET", "https://example.com/405", nil),
			code:         405,
			marshal404:   marshal404,
			marshal405:   marshal405,
			buf:          []byte(http.StatusText(http.StatusMethodNotAllowed)),
			expectedCode: 405,
			expectedBody: []byte(`{"code":405}`),
		},
		"status 405 with marshal error": {
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest("GET", "https://example.com/405", nil),
			code:         405,
			marshal404:   marshal404,
			marshal405:   marshalErr,
			buf:          []byte(http.StatusText(http.StatusMethodNotAllowed)),
			expectedCode: 405,
			expectedBody: nil,
		},
		"status not 404 and 405": {
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest("GET", "https://example.com/200", nil),
			code:         200,
			marshal404:   marshal404,
			marshal405:   marshal405,
			buf:          []byte(`{"code":200}`),
			expectedCode: 200,
			expectedBody: []byte(`{"code":200}`),
		},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			wwr := nrmux.WrapResponseWriter{ResponseWriter: v.w}
			wwr.SetCode(v.code)
			wwr.SetReq(v.r)
			wwr.SetMarshal404(v.marshal404)
			wwr.SetMarshal405(v.marshal405)

			wwr.Write(v.buf)

			require.Equal(t, v.expectedCode, v.w.Code)
			require.Equal(t, v.expectedBody, v.w.Body.Bytes())
			require.Equal(t, "application/json", v.w.Header().Get("Content-Type"))
		})
	}
}
