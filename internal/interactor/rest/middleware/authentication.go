package middleware

import (
	"context"
	"fmt"
	"go-playground/internal/interactor/rest/preauth"
	"go-playground/pkg/render"
	"net/http"

	"go.uber.org/zap"
)

const authHeader = "Authorization"

// Autheticator authenticate Authorization Header token via AuthenticationManager.
func Autheticator(logger *zap.Logger, manager preauth.AuthenticationManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := manager.Authenticate(r.Header.Get(authHeader))
			if err != nil {
				logger.Warn("No authenticated request had recieved.", zap.Error(err))
				render.Unauthorized(w, "You had failed to authenticate", r.URL.Path)
				return
			}
			ctx := context.WithValue(r.Context(), authCtxKey{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
		return fn
	}
}

// GetAuthenticatedUser retrive authenticated user from context.
// This variable useful for mocking behavior.
var GetAutenticatedUser = getAuthenticatedUserFromContext

// GetAuthenticatedUserFromContext retrive authenticated user information from context.
// If errro is not nil, this indicates request is not authenticated.
func getAuthenticatedUserFromContext(ctx context.Context) (*preauth.AuthenticatedUser, error) {
	u, ok := ctx.Value(authCtxKey{}).(*preauth.AuthenticatedUser)
	if !ok || u == nil {
		return nil, fmt.Errorf("user does not exist context")
	}
	return u, nil
}

// authCtxKey is context key storeed AuthenticateUser.
type authCtxKey struct{}
