package middleware

import (
	"context"
	"fmt"
	"go-playground/internal/transport_layer/rest/preauth"
	"go-playground/pkg/renderer"
	"net/http"

	"go.uber.org/zap"
)

const authHeader = "Authorization"

// Authenticator authenticate Authorization Header token via AuthenticationManager.
func Authenticator(logger *zap.Logger, manager preauth.AuthenticationManager, failure renderer.Failure) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := manager.Authenticate(r.Header.Get(authHeader))
			if err != nil {
				logger.Warn("No authenticated request had received.", zap.Error(err))
				failure.Response(w, r, http.StatusUnauthorized, "Request With No Authentication", "Request token was not found in your request header")
				return
			}
			ctx := context.WithValue(r.Context(), authCtxKey{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
		return fn
	}
}

// CurrentUser retrieve authenticated user from context.
// This variable useful for mocking behavior.
var CurrentUser = getAuthenticatedUserFromContext

// GetAuthenticatedUserFromContext retrieve authenticated user information from context.
// If error is not nil, this indicates request is not authenticated.
func getAuthenticatedUserFromContext(ctx context.Context) (*preauth.AuthenticatedUser, error) {
	u, ok := ctx.Value(authCtxKey{}).(*preauth.AuthenticatedUser)
	if !ok || u == nil {
		return nil, fmt.Errorf("user does not exist context")
	}
	return u, nil
}

// authCtxKey is context key stored AuthenticateUser.
type authCtxKey struct{}
