package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReplyHandler(t *testing.T) {
	tests := map[string]struct {
		w            *httptest.ResponseRecorder
		r            *http.Request
		expectedCode int
		expectedBody []byte
	}{
		"request to /reply/Ryan": {
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest("GET", "https://example.com/reply/Ryan", nil),
			expectedCode: 200,
			expectedBody: []byte(`{"message":"Hi Ryan"}` + "\n"),
		},
		"request to /reply/": {
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest("GET", "https://example.com/reply/", nil),
			expectedCode: 404,
			expectedBody: []byte(`404 page not found` + "\n"),
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
