package handler_test

import (
	"context"
	"go-playground/cmd/api/internal/transportlayer/rest/handler/v2"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := map[string]struct {
		setup   func(t *testing.T)
		wantErr bool
	}{
		"success": {
			setup: func(t *testing.T) {
				t.Setenv("DB_USER", "dummy_user")
				t.Setenv("DB_PASSWORD", "dummy_password")
				t.Setenv("DB_ADDRESS", "localhost:3306")
				t.Setenv("DB_NAME", "dummy")
				t.Setenv("AUTH_ISSUER_URL", "http://example.com")
			},
		},
		"failure: failed to create query db": {
			setup:   func(t *testing.T) { /* noop */ },
			wantErr: true,
		},
		"failure: failed to find issuer url": {
			setup: func(t *testing.T) {
				t.Setenv("DB_USER", "dummy_user")
				t.Setenv("DB_PASSWORD", "dummy_password")
				t.Setenv("DB_ADDRESS", "localhost:3306")
				t.Setenv("DB_NAME", "dummy")
			},
			wantErr: true,
		},
		"failure: failed to create auth middleware": {
			setup: func(t *testing.T) {
				t.Setenv("DB_USER", "dummy_user")
				t.Setenv("DB_PASSWORD", "dummy_password")
				t.Setenv("DB_ADDRESS", "localhost:3306")
				t.Setenv("DB_NAME", "dummy")
				t.Setenv("AUTH_ISSUER_URL", "")
			},
			wantErr: true,
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			v.setup(t)

			gotHandler, gotErr := handler.New(nil, os.LookupEnv)
			if v.wantErr {
				assert.Empty(t, gotHandler)
				assert.Error(t, gotErr)
			} else {
				assert.NotEmpty(t, gotHandler)
				assert.NoError(t, gotErr)
			}
		})
	}
}

func TestNew_ErrorHandlerFunc(t *testing.T) {
	t.Setenv("DB_USER", "dummy_user")
	t.Setenv("DB_PASSWORD", "dummy_password")
	t.Setenv("DB_ADDRESS", "localhost:3306")
	t.Setenv("DB_NAME", "dummy")
	t.Setenv("AUTH_ISSUER_URL", "http://example.com")
	hn, err := handler.New(nil, os.LookupEnv)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/tasks?limit=not_number", nil)

	hn.ServeHTTP(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"message":"Invalid format for parameter limit: error binding string parameter: strconv.ParseInt: parsing \"not_number\": invalid syntax"}`, w.Body.String())
}
