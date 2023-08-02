package middleware

import (
	"context"
	"fmt"
	"go-playground/pkg/lib/render"
	"go-playground/pkg/presentation/preauth"
	"net/http"

	"go.uber.org/zap"
)

const authHeader = "Authorization"

type authenticationMiddleWare struct {
	logger      *zap.Logger
	authManager preauth.AuthenticationManager
}

// NewAuthenticationMiddleWare init Middleware interface.
func NewAuthenticationMiddleWare(logger *zap.Logger, authenticationManager preauth.AuthenticationManager) MiddleWare {
	return &authenticationMiddleWare{logger, authenticationManager}
}

// Handle to store authenticated user info in context when user request is authenticated.
// If requests is not authenticated, return 401 staus code to client.
func (mid *authenticationMiddleWare) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := mid.authManager.Authenticate(r.Header.Get(authHeader))
		if err != nil {
			mid.logger.Warn("No authenticated request had recieved.", zap.String("path", r.URL.Path), zap.Error(err))
			render.Unauthorized(w, "You had failed to authenticate", r.URL.Path)
		} else {
			ctx := context.WithValue(r.Context(), authCtxKey{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
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
