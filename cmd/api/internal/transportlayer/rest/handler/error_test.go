package handler_test

import (
	"errors"
	"go-playground/cmd/api/internal/transportlayer/rest/handler"
	"go-playground/pkg/errorx"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorHandlerFunc(t *testing.T) {
	tests := map[string]struct {
		w            *httptest.ResponseRecorder
		r            *http.Request
		fn           handler.ErrorHandlerFunc
		expectedCode int
		expectedBody string
	}{
		"success": {
			w: httptest.NewRecorder(),
			r: httptest.NewRequest("", "/", nil),
			fn: handler.ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
				_, err := w.Write([]byte(`{"msg":"ok"}`))
				return err
			}),
			expectedCode: 200,
			expectedBody: `{"msg":"ok"}`,
		},
		"unhandled error": {
			w: httptest.NewRecorder(),
			r: httptest.NewRequest("", "/", nil),
			fn: handler.ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
				return errors.New("unhandled")
			}),
			expectedCode: 500,
			expectedBody: `{"message":"unhandled"}`,
		},
		"handled error": {
			w: httptest.NewRecorder(),
			r: httptest.NewRequest("", "/", nil),
			fn: handler.ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
				return errorx.NewError("handled error", errorx.WithStatus(400))
			}),
			expectedCode: 400,
			expectedBody: `{"message":"handled error"}`,
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			v.fn.ServeHTTP(v.w, v.r)

			assert.Equal(t, v.expectedCode, v.w.Code)
			assert.JSONEq(t, v.expectedBody, v.w.Body.String())
		})
	}
}
