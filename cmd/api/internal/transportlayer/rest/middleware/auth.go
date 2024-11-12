package middleware

import (
	"errors"
	"fmt"
	"go-playground/cmd/api/internal/transportlayer/rest"
	"go-playground/pkg/ctxhelper"
	"net/http"
	"strings"

	jwt "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/newrelic/go-agent/v3/newrelic"
)

var errUnexpectedClaim = errors.New("unexpected jwt claim")

// Auth check access token.
type Auth struct {
	jwsValidator *validator.Validator
	skipRoutes   map[string]struct{}
}

// NewAuth creates [Auth].
func NewAuth(
	jwksProvider *jwks.CachingProvider,
	audience []string,
	opts ...authOption,
) (*Auth, error) {
	v, err := validator.New(
		jwksProvider.KeyFunc,
		validator.RS256,
		jwksProvider.IssuerURL.String(),
		audience,
	)
	if err != nil {
		return nil, fmt.Errorf("new validator for auth: %w", err)
	}
	auth := Auth{jwsValidator: v, skipRoutes: make(map[string]struct{})}
	for _, fn := range opts {
		fn(&auth)
	}
	return &auth, nil
}

// CheckAccessToken checks access token attached Authorization header.
func (a *Auth) CheckAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txn := newrelic.FromContext(r.Context())
		segment := txn.StartSegment("middleware/Auth/CheckAccessToken")
		defer segment.End()

		_, skip := a.skipRoutes[r.Pattern]
		if r.Method == http.MethodOptions || skip {
			next.ServeHTTP(w, r)
			return
		}

		jws, ok := a.extractJWS(r.Header.Get("Authorization"))
		if !ok {
			rest.Err(w, "missing access token", http.StatusBadRequest)
			return
		}
		validated, err := a.jwsValidator.ValidateToken(r.Context(), jws)
		if err != nil {
			if errors.Is(err, jwt.ErrJWTInvalid) { // TODO: correct handling.
				rest.Err(w, "invalid access token", http.StatusUnauthorized)
				return
			}
			rest.Err(w, "unexpected error when checking access token", http.StatusInternalServerError)
			txn.NoticeError(err)
			return
		}
		claim, ok := validated.(*validator.ValidatedClaims)
		if !ok {
			rest.Err(w, "unexpected claim", http.StatusInternalServerError)
			txn.NoticeError(errUnexpectedClaim)
			return
		}
		ctx := ctxhelper.WithSubject(r.Context(), claim.RegisteredClaims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Auth) extractJWS(header string) (string, bool) {
	if header == "" {
		return "", false
	}
	parts := strings.Fields(header)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", false
	}
	return parts[1], true
}
