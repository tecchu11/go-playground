package middleware

import (
	"context"
	"errors"
	"fmt"
	"go-playground/pkg/renderer"
	"net/http"

	"go.uber.org/zap"
)

// UserRole replesents user role to perform to request.
type UserRole int

const (
	// Admin is super user
	Admin UserRole = iota + 1
	// User is role for general purpose.
	User
)

var (
	roleStrings = map[UserRole]string{Admin: "Admin", User: "User"}
	stringRoles = map[string]UserRole{"Admin": Admin, "User": User}
)

// String convert UserRole to string type.
// UserRole.Admin to "Admin", UserRole.User to "User"
func (ur UserRole) String() string {
	return roleStrings[ur]
}

// UserRoleString lookup UserRole by given string.
func UserRoleString(v string) (UserRole, error) {
	ur, ok := stringRoles[v]
	if !ok {
		return 0, fmt.Errorf("role %s does not defined", v)
	}
	return ur, nil
}

type (
	// AuthUser replesents pre authenticated user.
	AuthUser struct {
		Name string
		Role UserRole
	}
	// PreAuthenticatedUsers replesents pre authenticated user map.
	// Please set key as token that authenticated user's.
	PreAuthenticatedUsers map[string]AuthUser
	authCtxKey            struct{}
)

const (
	authHeader = "Authorization"
)

var (
	NoUser    = AuthUser{}
	ErrNoUser = errors.New("there is no auth use")
)

// Auth middleware for authentication and authorization by pre-authenticated token and given permitRoles.
func (u PreAuthenticatedUsers) Auth(logger *zap.Logger, failure renderer.JSON, permitRoles map[UserRole]struct{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := u[r.Header.Get(authHeader)]
			if !ok {
				failure.Failure(w, r, http.StatusUnauthorized, "Request With No Authentication", "Request token was not found in your request header")
				return
			}
			_, ok = permitRoles[user.Role]
			if !ok {
				failure.Failure(w, r, http.StatusForbidden, "Request With No Authorization", fmt.Sprintf("Your role(%s) was not performing to action", user.Role.String()))
				return
			}
			next.ServeHTTP(w, r.WithContext(user.Set(r.Context())))
		})
		return fn
	}
}

// CurrentUser retreive from AuthUser from context.
func CurrentUser(ctx context.Context) (AuthUser, error) {
	user, ok := ctx.Value(authCtxKey{}).(AuthUser)
	if !ok || user == NoUser {
		return NoUser, ErrNoUser
	}
	return user, nil
}

// Set AuthUser to request context.
func (u AuthUser) Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, authCtxKey{}, u)
}
