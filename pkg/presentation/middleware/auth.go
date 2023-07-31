package middleware

import (
	"context"
	"fmt"
	"go-playground/pkg/lib/render"
	"go-playground/pkg/presentation/auth"
	"net/http"

	"go.uber.org/zap"
)

const authHeader = "Authorization"

type authMiddleWare struct {
	logger      *zap.Logger
	authManager auth.AuthenticationManager
}

// NewAuthMiddleWare init Middleware interface.
func NewAuthMiddleWare(logger *zap.Logger, authenticationManager auth.AuthenticationManager) MiddleWare {
	return &authMiddleWare{logger, authenticationManager}
}

// Handle to store authenticated user info in context when user request is authenticated.
// If requests is not authenticated, return 401 staus code to client.
func (mid *authMiddleWare) Handle(next http.Handler) http.Handler {
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

// GetAuthenticatedUserFromContext retrive authenticated user information from context.
// If errro is not nil, this indicates request is not authenticated.
func GetAuthenticatedUserFromContext(ctx context.Context) (*auth.AuthenticatedUser, error) {
	u, ok := ctx.Value(authCtxKey{}).(*auth.AuthenticatedUser)
	if !ok || u == nil {
		return nil, fmt.Errorf("user does not exist context")
	}
	return u, nil
}

// authCtxKey is context key storeed AuthenticateUser.
type authCtxKey struct{}

// AuthCtxManager implemented contextutil.ContextManger
type AuthCtxManager struct{}

// Get *AuthenticatedUser from context.
// If errro is not nil, this indicates request is not authenticated.
func (manager *AuthCtxManager) Get(ctx context.Context) (*auth.AuthenticatedUser, error) {
	u, ok := ctx.Value(authCtxKey{}).(*auth.AuthenticatedUser)
	if !ok || u == nil {
		return nil, fmt.Errorf("user does not exist context")
	}
	return u, nil
}
