package handler_test

import (
	"context"
	"errors"
	"go-playground/cmd/api/internal/transportlayer/rest/handler"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type pingerFunc func(context.Context) error

func (ping pingerFunc) PingContext(ctx context.Context) error {
	return ping(ctx)
}

func TestHealthCheck(t *testing.T) {
	tests := map[string]struct {
		pinger pingerFunc
		want   struct {
			code int
			body string
		}
	}{
		"success": {
			pinger: func(ctx context.Context) error { return nil },
			want: struct {
				code int
				body string
			}{code: 200, body: `{"msg":"ok"}`},
		},
		"fail ping": {
			pinger: func(ctx context.Context) error { return errors.New("failed to ping") },
			want: struct {
				code int
				body string
			}{code: 500, body: "null"},
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/health", nil)

			handler.HealthCheck(v.pinger).ServeHTTP(w, r)

			assert.Equal(t, v.want.code, w.Code)
			assert.JSONEq(t, v.want.body, w.Body.String())
		})
	}
}
