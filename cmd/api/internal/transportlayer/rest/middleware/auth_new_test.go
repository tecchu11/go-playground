package middleware_test

import (
	"go-playground/cmd/api/internal/transportlayer/rest/middleware"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewCheckAccessToken(t *testing.T) {
	type input struct {
		cfg middleware.CheckAccessTokenConfig
	}
	type want struct {
		handlerNil bool
		errNil     bool
	}
	tests := map[string]struct {
		input input
		want  want
	}{
		"success": {
			input: input{
				cfg: middleware.CheckAccessTokenConfig{
					IssuerURL:     mustParseURL(t, "https://example.com"),
					Audiences:     []string{"api"},
					ExclusionURLs: []string{"/health"},
					HTTPClient:    &http.Client{},
					CacheTTL:      time.Hour,
					Logger:        slog.New(slog.NewTextHandler(os.Stderr, nil)),
				},
			},
			want: want{
				handlerNil: false,
				errNil:     true,
			},
		},
		"success with CacheTTL=0(fallback to library default value)": {
			input: input{
				cfg: middleware.CheckAccessTokenConfig{
					IssuerURL:     mustParseURL(t, "https://example.com"),
					Audiences:     []string{"api"},
					ExclusionURLs: []string{"/health"},
					HTTPClient:    &http.Client{},
					CacheTTL:      0,
					Logger:        slog.New(slog.NewTextHandler(os.Stderr, nil)),
				},
			},
			want: want{
				handlerNil: false,
				errNil:     true,
			},
		},
		"success with empty ExclusionURLs": {
			input: input{
				cfg: middleware.CheckAccessTokenConfig{
					IssuerURL:     mustParseURL(t, "https://example.com"),
					Audiences:     []string{"api"},
					ExclusionURLs: []string{},
					HTTPClient:    &http.Client{},
					CacheTTL:      time.Hour,
					Logger:        slog.New(slog.NewTextHandler(os.Stderr, nil)),
				},
			},
			want: want{
				handlerNil: false,
				errNil:     true,
			},
		},
		"error: IssuerURL is nil": {
			input: input{
				cfg: middleware.CheckAccessTokenConfig{
					IssuerURL:     nil,
					Audiences:     []string{"api"},
					ExclusionURLs: []string{},
					HTTPClient:    &http.Client{},
					CacheTTL:      time.Hour,
					Logger:        slog.New(slog.NewTextHandler(os.Stderr, nil)),
				},
			},
			want: want{
				handlerNil: true,
				errNil:     false,
			},
		},
		"error: HTTPClient is nil": {
			input: input{
				cfg: middleware.CheckAccessTokenConfig{
					IssuerURL:     mustParseURL(t, "https://example.com"),
					Audiences:     []string{"api"},
					ExclusionURLs: []string{},
					HTTPClient:    nil,
					CacheTTL:      time.Hour,
					Logger:        slog.New(slog.NewTextHandler(os.Stderr, nil)),
				},
			},
			want: want{
				handlerNil: true,
				errNil:     false,
			},
		},
		"error: CacheTTL < 0": {
			input: input{
				cfg: middleware.CheckAccessTokenConfig{
					IssuerURL:     mustParseURL(t, "https://example.com"),
					Audiences:     []string{"api"},
					ExclusionURLs: []string{},
					HTTPClient:    &http.Client{},
					CacheTTL:      -1 * time.Second,
					Logger:        slog.New(slog.NewTextHandler(os.Stderr, nil)),
				},
			},
			want: want{
				handlerNil: true,
				errNil:     false,
			},
		},
		"error: Audiences is empty": {
			input: input{
				cfg: middleware.CheckAccessTokenConfig{
					IssuerURL:     mustParseURL(t, "https://example.com"),
					Audiences:     []string{},
					ExclusionURLs: []string{},
					HTTPClient:    &http.Client{},
					CacheTTL:      time.Hour,
					Logger:        slog.New(slog.NewTextHandler(os.Stderr, nil)),
				},
			},
			want: want{
				handlerNil: true,
				errNil:     false,
			},
		},
		"error: Logger is nil": {
			input: input{
				cfg: middleware.CheckAccessTokenConfig{
					IssuerURL:     mustParseURL(t, "https://example.com"),
					Audiences:     []string{"audience"},
					ExclusionURLs: []string{},
					HTTPClient:    &http.Client{},
					CacheTTL:      time.Hour,
				},
			},
			want: want{
				handlerNil: true,
				errNil:     false,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			handler, err := middleware.NewCheckAccessToken(tt.input.cfg)

			if tt.want.errNil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			if tt.want.handlerNil {
				require.Nil(t, handler)
			} else {
				require.NotNil(t, handler)
			}
		})
	}
}

func mustParseURL(t *testing.T, urlStr string) *url.URL {
	u, err := url.Parse(urlStr)
	require.NoError(t, err)
	return u
}
