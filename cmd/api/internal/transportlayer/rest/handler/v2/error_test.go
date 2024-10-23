package handler_test

import (
	"go-playground/cmd/api/internal/transportlayer/rest/handler/v2"
	"go-playground/pkg/errorx"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorHandlerFunc(t *testing.T) {
	type input struct {
		w       *httptest.ResponseRecorder
		r       *http.Request
		handler func(http.ResponseWriter, *http.Request) error
	}
	type want struct {
		status int
		body   string
	}
	tests := map[string]struct {
		input input
		want  want
	}{
		"success": {
			input: input{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "http://example.com", nil),
				handler: func(w http.ResponseWriter, r *http.Request) error {
					_, err := w.Write([]byte(`{"message":"ok"}`))
					return err
				},
			},
			want: want{
				status: http.StatusOK,
				body:   `{"message":"ok"}`,
			},
		},
		"handler returns app error(unexpected)": {
			input: input{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "http://example.com", nil),
				handler: func(w http.ResponseWriter, r *http.Request) error {
					return errorx.NewError("app unexpected error", errorx.WithStatus(http.StatusInternalServerError))
				},
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   `{"message":"app unexpected error"}`,
			},
		},
		"handler returns app error(expected)": {
			input: input{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "http://example.com", nil),
				handler: func(w http.ResponseWriter, r *http.Request) error {
					return errorx.NewWarn("app expected error", errorx.WithStatus(http.StatusBadRequest))
				},
			},
			want: want{
				status: http.StatusBadRequest,
				body:   `{"message":"app expected error"}`,
			},
		},
		"handler returns unhandled error": {
			input: input{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "http://example.com", nil),
				handler: func(w http.ResponseWriter, r *http.Request) error {
					return io.ErrUnexpectedEOF
				},
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   `{"message":"unexpected EOF"}`,
			},
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			handler.ErrorHandlerFunc(v.input.w, v.input.r, v.input.handler)

			assert.Equal(t, v.want.status, v.input.w.Code)
			assert.JSONEq(t, v.want.body, v.input.w.Body.String())
		})
	}
}
