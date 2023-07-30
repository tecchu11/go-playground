package auth_test

import (
	"go-playground/pkg/presentation/auth"
	"testing"
)

func TestString(t *testing.T) {
	cases := []struct {
		name     string
		role     auth.Role
		expected string
	}{
		{
			name:     "case ADMIN role",
			role:     auth.ADMIN,
			expected: "ADMIN",
		},
		{
			name:     "case USER role",
			role:     auth.USER,
			expected: "USER",
		},
		{
			name:     "case UNDIFINED role",
			role:     auth.UNDIFINED,
			expected: "",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.role.String()
			if c.expected != actual {
				t.Errorf("RoleFrom result is unmatched. Expected is %v but actual is %v", c.expected, actual)
			}
		})
	}
}

func TestRoleFrom(t *testing.T) {
	cases := []struct {
		name     string
		literal  string
		expected auth.Role
	}{
		{
			name:     "case of ADMIN",
			literal:  "ADMIN",
			expected: auth.ADMIN,
		},
		{
			name:     "case of USER",
			literal:  "USER",
			expected: auth.USER,
		},
		{
			name:     "case of invalid",
			literal:  "Invalid role",
			expected: auth.UNDIFINED,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := auth.RoleFrom(c.literal)
			if c.expected != actual {
				t.Errorf("RoleFrom result is unmatched. Expected is %v but actual is %v", c.expected, actual)
			}

		})
	}
}
