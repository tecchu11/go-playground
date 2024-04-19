package router_test

import (
	"encoding/json"
	"errors"
	"go-playground/pkg/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInterceptWriter_Write(t *testing.T) {
	tests := map[string]struct {
		code       int
		w          *httptest.ResponseRecorder
		r          *http.Request
		marshal404 router.ErrMarshalFunc
		marshal405 router.ErrMarshalFunc
		expected   []byte
	}{
		"with default marshal on 404": {
			code:       404,
			w:          httptest.NewRecorder(),
			r:          httptest.NewRequest("GET", "https://example.com/404", nil),
			marshal404: router.DefaultErrMarshal,
			marshal405: router.DefaultErrMarshal,
		},
		"with custom marshal on 404": {
			code: 404,
			w:    httptest.NewRecorder(),
			r:    httptest.NewRequest("GET", "https://example.com/404", nil),
			marshal404: func(_ *http.Request) ([]byte, error) {
				return json.Marshal(map[string]string{"message": "404 error"})
			},
			marshal405: router.DefaultErrMarshal,
			expected:   []byte(`{"message":"404 error"}`),
		},
		"with invalid custom marshal on 404": {
			code: 404,
			w:    httptest.NewRecorder(),
			r:    httptest.NewRequest("GET", "https://example.com/404", nil),
			marshal404: func(_ *http.Request) ([]byte, error) {
				return nil, errors.New("marshal error")
			},
			marshal405: router.DefaultErrMarshal,
		},
		"with default marshal on 405": {
			code:       405,
			w:          httptest.NewRecorder(),
			r:          httptest.NewRequest("GET", "https://example.com/405", nil),
			marshal404: router.DefaultErrMarshal,
			marshal405: router.DefaultErrMarshal,
		},
		"with custom marshal on 405": {
			code:       405,
			w:          httptest.NewRecorder(),
			r:          httptest.NewRequest("GET", "https://example.com/405", nil),
			marshal404: router.DefaultErrMarshal,
			marshal405: func(_ *http.Request) ([]byte, error) {
				return json.Marshal(map[string]string{"message": "405 error"})
			},
			expected: []byte(`{"message":"405 error"}`),
		},
		"with invalid custom marshal on 405": {
			code:       405,
			w:          httptest.NewRecorder(),
			r:          httptest.NewRequest("GET", "https://example.com/405", nil),
			marshal404: router.DefaultErrMarshal,
			marshal405: func(_ *http.Request) ([]byte, error) {
				return nil, errors.New("marshal error")
			},
		},
		"neither 404 nor 405": {
			code:       200,
			w:          httptest.NewRecorder(),
			r:          httptest.NewRequest("GET", "https://example.com/200", nil),
			marshal404: router.DefaultErrMarshal,
			marshal405: router.DefaultErrMarshal,
			expected:   []byte("Response"),
		},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			iw := router.InterceptWriter{
				ResponseWriter: v.w,
				Req:            v.r,
				Marshal404:     v.marshal404,
				Marshal405:     v.marshal405,
			}
			iw.WriteHeader(v.code)
			_, err := iw.Write([]byte("Response"))
			require.NoError(t, err)
			require.Equal(t, v.code, v.w.Code)
			require.Equal(t, v.expected, v.w.Body.Bytes())
		})
	}
}

func TestInterceptWriter_WriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "https://example.com/notfound", nil)

	iw := router.InterceptWriter{
		ResponseWriter: w,
		Req:            r,
		Marshal404:     router.DefaultErrMarshal,
		Marshal405:     router.DefaultErrMarshal,
	}
	iw.WriteHeader(404)
	require.Equal(t, 200, w.Code)
	require.Equal(t, 404, iw.Code)
}
