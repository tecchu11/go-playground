package auth_test

import (
	"context"
	"errors"
	"go-playground/pkg/presentation/auth"
	"reflect"
	"testing"
)

func TestGetAuthUser(t *testing.T) {
	cases := []struct {
		name         string
		ctx          context.Context
		expectedUser *auth.AuthenticatedUser
		expectedErr  error
	}{
		{
			name: "case of successful to retrive user from context",
			ctx: context.WithValue(
				context.Background(),
				"authUser",
				&auth.AuthenticatedUser{Name: "test-user", Role: auth.ADMIN},
			),
			expectedUser: &auth.AuthenticatedUser{Name: "test-user", Role: auth.ADMIN},
		},
		{
			name:        "case of failuer to retive user from context",
			ctx:         context.Background(),
			expectedErr: auth.ErrNoUser,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			acutualUser, actualErr := auth.GetAuthUser(c.ctx)
			if !reflect.DeepEqual(acutualUser, c.expectedUser) {
				t.Errorf("user retrived from context does not match between actual %v and expected %v", acutualUser, c.expectedUser)
			}
			if !errors.Is(actualErr, c.expectedErr) {
				t.Errorf("error does not match between (%v) and (%v)", actualErr, c.expectedErr)
			}
		})
	}
}

func TestAuthenticatedUser_SetContext(t *testing.T) {
	user := &auth.AuthenticatedUser{
		Name: "test-user",
		Role: auth.ADMIN,
	}
	expected := context.WithValue(
		context.Background(),
		"authUser", user,
	)
	actual := user.SetContext(context.Background())
	if ok := reflect.DeepEqual(expected, actual); !ok {
		t.Errorf("SetContext is unexpected because expecte is %v but actual is %v", expected, actual)
	}
}
