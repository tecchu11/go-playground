package middleware

import (
	"go-playground/pkg/presentation/auth"
	"go-playground/pkg/presentation/handler"
	"go-playground/pkg/presentation/model"
	"net/http"

	"go.uber.org/zap"
)

const authHeader = "Authorization"

type authMiddleWare struct {
	logger      *zap.Logger
	authManager *auth.AuthenticationManager
}

func NewAuthMiddleWare(logger *zap.Logger, authManager *auth.AuthenticationManager) MiddleWare {
	return &authMiddleWare{logger: logger, authManager: authManager}
}

func (mid *authMiddleWare) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := mid.authManager.Authenticate(r.Header.Get(authHeader))
		if err != nil {
			mid.logger.Warn("No authenticated request had recieved.", zap.String("path", r.URL.Path), zap.Error(err))
			res := model.NewProblemDetail("You had failed to authenticate", r.URL.Path, http.StatusUnauthorized)
			handler.JsonResponse(w, http.StatusUnauthorized, res)
		} else {
			ctx := user.SetContext(r.Context())
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
