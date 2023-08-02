package auth_test

import (
	"go-playground/pkg/presentation/auth"
	"testing"
)

func Test_Authorize(t *testing.T) {
	tests := map[string]struct {
		authorizedList []auth.Role
		inputRole      auth.Role
		expectedErr    bool
	}{
		"test permited admin only":                               {[]auth.Role{auth.ADMIN}, auth.ADMIN, false},
		"test permited admin only and then passed user":          {[]auth.Role{auth.ADMIN}, auth.USER, true},
		"test permited user only":                                {[]auth.Role{auth.USER}, auth.USER, false},
		"test permited user only and then passed admin":          {[]auth.Role{auth.USER}, auth.ADMIN, true},
		"test permited user and admin and then passed user":      {[]auth.Role{auth.USER, auth.ADMIN}, auth.USER, false},
		"test permited user and admin and then passed admin":     {[]auth.Role{auth.USER, auth.ADMIN}, auth.ADMIN, false},
		"test permited admin only and then passed undifined":     {[]auth.Role{auth.ADMIN}, auth.UNDIFINED, true},
		"test permited user only and then passed undifined":      {[]auth.Role{auth.USER}, auth.UNDIFINED, true},
		"test permited user and admin and then passed undifined": {[]auth.Role{auth.USER, auth.ADMIN}, auth.UNDIFINED, true},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			var list auth.AuthorizedList = v.authorizedList
			err := list.Authorize(v.inputRole)
			if v.expectedErr && err == nil {
				t.Error("err is nil but expect err")
			}
			if !v.expectedErr && err != nil {
				t.Errorf("no err is expected but acutal err (%s) happend", err)
			}
		})
	}
}
