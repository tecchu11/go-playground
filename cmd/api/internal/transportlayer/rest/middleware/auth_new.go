package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v3"
	"github.com/auth0/go-jwt-middleware/v3/jwks"
	"github.com/auth0/go-jwt-middleware/v3/validator"
)

// CheckAccessTokenConfig holds configuration for JWT validation middleware.
type CheckAccessTokenConfig struct {
	IssuerURL     *url.URL
	Audiences     []string
	ExclusionURLs []string
	HTTPClient    *http.Client
	CacheTTL      time.Duration
	Logger        jwtmiddleware.Logger
}

// NewCheckAccessToken creates JWT validation middleware using the official v3 middleware.
// It initializes the JWKS provider and validator internally, returning an HTTP middleware
// that validates JWT tokens and skips the specified URLs.
func NewCheckAccessToken(cfg CheckAccessTokenConfig) (func(http.Handler) http.Handler, error) {
	jwksProvider, err := jwks.NewCachingProvider(
		jwks.WithIssuerURL(cfg.IssuerURL),
		jwks.WithCacheTTL(cfg.CacheTTL),
		jwks.WithCustomClient(cfg.HTTPClient),
	)
	if err != nil {
		return nil, fmt.Errorf("new jwks provider: %w", err)
	}

	v, err := validator.New(
		validator.WithKeyFunc(jwksProvider.KeyFunc),
		validator.WithAlgorithm(validator.RS256),
		validator.WithIssuer(cfg.IssuerURL.String()),
		validator.WithAudiences(cfg.Audiences),
	)
	if err != nil {
		return nil, fmt.Errorf("new validator: %w", err)
	}

	var opts []jwtmiddleware.Option
	opts = append(opts, jwtmiddleware.WithValidator(v), jwtmiddleware.WithLogger(cfg.Logger))
	if len(cfg.ExclusionURLs) > 0 {
		opts = append(opts, jwtmiddleware.WithExclusionUrls(cfg.ExclusionURLs))
	}

	middleware, err := jwtmiddleware.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("new jwt middleware: %w", err)
	}

	return middleware.CheckJWT, nil
}
