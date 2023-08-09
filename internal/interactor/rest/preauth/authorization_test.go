package preauth_test

import (
	preauth2 "go-playground/internal/interactor/rest/preauth"
	"testing"
)

func Test_Authorize(t *testing.T) {
	tests := map[string]struct {
		authorizedList []preauth2.Role
		inputRole      preauth2.Role
		expectedErr    bool
	}{
		"test permited admin only":                               {[]preauth2.Role{preauth2.ADMIN}, preauth2.ADMIN, false},
		"test permited admin only and then passed user":          {[]preauth2.Role{preauth2.ADMIN}, preauth2.USER, true},
		"test permited user only":                                {[]preauth2.Role{preauth2.USER}, preauth2.USER, false},
		"test permited user only and then passed admin":          {[]preauth2.Role{preauth2.USER}, preauth2.ADMIN, true},
		"test permited user and admin and then passed user":      {[]preauth2.Role{preauth2.USER, preauth2.ADMIN}, preauth2.USER, false},
		"test permited user and admin and then passed admin":     {[]preauth2.Role{preauth2.USER, preauth2.ADMIN}, preauth2.ADMIN, false},
		"test permited admin only and then passed undifined":     {[]preauth2.Role{preauth2.ADMIN}, preauth2.UNDIFINED, true},
		"test permited user only and then passed undifined":      {[]preauth2.Role{preauth2.USER}, preauth2.UNDIFINED, true},
		"test permited user and admin and then passed undifined": {[]preauth2.Role{preauth2.USER, preauth2.ADMIN}, preauth2.UNDIFINED, true},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			authorizationManager := preauth2.NewAuthorizationManager(v.authorizedList)
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
