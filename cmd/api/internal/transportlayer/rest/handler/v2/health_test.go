package handler_test

import (
	"context"
	"errors"
	"go-playground/cmd/api/internal/transportlayer/rest/handler/v2"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	type input struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	type want struct {
		status int
		body   string
	}
	tests := map[string]struct {
		input input
		setup func() *handler.HealthHandler
		want  want
	}{
		"success": {
			input: input{
				w: httptest.NewRecorder(),
				r: httptest.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/health", nil),
			},
			setup: func() *handler.HealthHandler {
				mck := new(MockPinger)
				mck.On("PingContext", context.Background()).Return(nil)
				return &handler.HealthHandler{Pinger: mck}
			},
			want: want{
				status: http.StatusOK,
				body:   `{"message":"ok"}`,
			},
		},
		"failed to ping": {
			input: input{
				w: httptest.NewRecorder(),
				r: httptest.NewRequestWithContext(context.Background(), http.MethodGet, "http://example.com/health", nil),
			},
			setup: func() *handler.HealthHandler {
				mck := new(MockPinger)
				mck.On("PingContext", context.Background()).Return(errors.New("failed to ping"))
				return &handler.HealthHandler{Pinger: mck}
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   `{"message":"internal server error"}`,
			},
		},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			hn := v.setup()

			hn.HealthCheck(v.input.w, v.input.r)

			assert.Equal(t, v.want.status, v.input.w.Code)
			assert.JSONEq(t, v.want.body, v.input.w.Body.String())
		})
	}
}
