package middleware_test

import (
	"context"
	"encoding/json"
	"go-playground/internal/transport_layer/rest/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserRole_String(t *testing.T) {
	tests := map[string]struct {
		in   middleware.UserRole
		want string
	}{
		"Admin role to string": {in: middleware.Admin, want: "Admin"},
		"User role to string":  {in: middleware.User, want: "User"},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			got := v.in.String()
			assert.Equal(t, v.want, got)
		})
	}
}

func TestUserRoleString(t *testing.T) {
	tests := map[string]struct {
		in        string
		expectErr bool
		want      middleware.UserRole
	}{
		"Admin to UserRole.Admin": {in: "Admin", expectErr: false, want: middleware.Admin},
		"User to UserRole.User":   {in: "User", expectErr: false, want: middleware.User},
		"Invalid to error":        {in: "Invalid", expectErr: true},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			got, gotErr := middleware.UserRoleString(v.in)
			if v.expectErr {
				assert.Error(t, gotErr)
				return
			}
			assert.Equal(t, v.want, got)
			assert.NoError(t, gotErr)
		})
	}
}

func TestPreAuthenticatedUsers_Auth_AdminOnly(t *testing.T) {
	tests := map[string]struct {
		inToken  string
		wantCode int
		wantBody map[string]string
	}{
		"response code is 200": {inToken: "admin-token", wantCode: 200, wantBody: map[string]string{"status": "ok"}},
		"response code is 403": {inToken: "user-token", wantCode: 403, wantBody: map[string]string{"detail": "Your role(User) was not performing to action", "title": "Request With No Authorization"}},
		"response code is 401": {inToken: "invalid-token", wantCode: 401, wantBody: map[string]string{"detail": "Request token was not found in your request header", "title": "Request With No Authentication"}},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/foos", nil)
			r.Header.Add("Authorization", v.inToken)
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(200)
				_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
			})
			users := middleware.PreAuthenticatedUsers(map[string]middleware.AuthUser{"user-token": {Name: "tecchu", Role: middleware.User}, "admin-token": {Name: "tecchu", Role: middleware.Admin}})

			adminOnly := users.Auth(&mockJSON{}, map[middleware.UserRole]struct{}{middleware.Admin: {}})
			adminOnly(nextHandler).ServeHTTP(w, r)

			var gotBody map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &gotBody)
			assert.NoError(t, err)
			assert.Equal(t, v.wantCode, w.Code)
			assert.Equal(t, v.wantBody, gotBody)
		})
	}
}

func TestCurrentUser(t *testing.T) {
	tests := map[string]struct {
		inCtx     context.Context
		expectErr bool
		want      middleware.AuthUser
	}{
		"success to retrieve user from context": {
			inCtx:     context.WithValue(context.Background(), middleware.AuthCtxKey, middleware.AuthUser{Name: "tecchu", Role: middleware.Admin}),
			expectErr: false,
			want:      middleware.AuthUser{Name: "tecchu", Role: middleware.Admin},
		},
		"error when no current user": {inCtx: context.Background(), expectErr: true},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			got, gotErr := middleware.CurrentUser(v.inCtx)
			if v.expectErr {
				assert.Error(t, gotErr)
				return
			}
			assert.Equal(t, v.want, got)
			assert.NoError(t, gotErr)
		})
	}
}

func TestAuthUser_Set(t *testing.T) {
	tests := map[string]struct {
		inUser  middleware.AuthUser
		inCtx   context.Context
		wantCtx context.Context
	}{
		"success to set user": {
			inUser:  middleware.AuthUser{Name: "tecchu", Role: middleware.Admin},
			inCtx:   context.Background(),
			wantCtx: context.WithValue(context.Background(), middleware.AuthCtxKey, middleware.AuthUser{Name: "tecchu", Role: middleware.Admin}),
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			got := v.inUser.Set(v.inCtx)
			assert.Equal(t, v.wantCtx, got)
		})
	}
}
