package middleware

import (
	"go-playground/pkg/lib/render"
	"go-playground/pkg/presentation/preauth"
	"net/http"

	"go.uber.org/zap"
)

type authorizationMiddleWare struct {
	logger         *zap.Logger
	authorizedList preauth.AuthorizedList
}

func NewAuthorizationMiddleWare(logger *zap.Logger, permitedRoles []preauth.Role) MiddleWare {
	var list preauth.AuthorizedList = permitedRoles
	return &authorizationMiddleWare{logger, list}
}

func (mid *authorizationMiddleWare) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := GetAutenticatedUser(r.Context())
		if err != nil {
			render.Unauthorized(w, "you are not authenticated because you requested with invalid token", r.URL.Path)
			return
		}
		if err := mid.authorizedList.Authorize(user.Role); err != nil {
			render.Forbidden(w, "you have no role to perform action", r.URL.Path)
			return
		}
		next.ServeHTTP(w, r)
	})
}
