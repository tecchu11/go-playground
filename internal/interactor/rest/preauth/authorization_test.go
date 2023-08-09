package preauth_test

import (
	"go-playground/internal/interactor/rest/preauth"
	"testing"
)

func Test_Authorize(t *testing.T) {
	tests := map[string]struct {
		authorizedList []preauth.Role
		inputRole      preauth.Role
		expectedErr    bool
	}{
		"test permited admin only":                               {[]preauth.Role{preauth.ADMIN}, preauth.ADMIN, false},
		"test permited admin only and then passed user":          {[]preauth.Role{preauth.ADMIN}, preauth.USER, true},
		"test permited user only":                                {[]preauth.Role{preauth.USER}, preauth.USER, false},
		"test permited user only and then passed admin":          {[]preauth.Role{preauth.USER}, preauth.ADMIN, true},
		"test permited user and admin and then passed user":      {[]preauth.Role{preauth.USER, preauth.ADMIN}, preauth.USER, false},
		"test permited user and admin and then passed admin":     {[]preauth.Role{preauth.USER, preauth.ADMIN}, preauth.ADMIN, false},
		"test permited admin only and then passed undifined":     {[]preauth.Role{preauth.ADMIN}, preauth.UNDIFINED, true},
		"test permited user only and then passed undifined":      {[]preauth.Role{preauth.USER}, preauth.UNDIFINED, true},
		"test permited user and admin and then passed undifined": {[]preauth.Role{preauth.USER, preauth.ADMIN}, preauth.UNDIFINED, true},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			authorizationManager := preauth.NewAuthorizationManager(v.authorizedList)
			err := authorizationManager.Authorize(v.inputRole)
			if v.expectedErr && err == nil {
				t.Error("err is nil but expect err")
			}
			if !v.expectedErr && err != nil {
				t.Errorf("no err is expected but acutal err (%s) happend", err)
			}
		})
	}
}
