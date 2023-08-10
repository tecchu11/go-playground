package preauth_test

import (
	"go-playground/internal/transport_layer/rest/preauth"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	cases := []struct {
		name     string
		role     preauth.Role
		expected string
	}{
		{name: "case ADMIN role", role: preauth.ADMIN, expected: "ADMIN"},
		{name: "case USER role", role: preauth.USER, expected: "USER"},
		{name: "case UNDEFINED role", role: preauth.UNDEFINED, expected: ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.role.String()
			assert.Equal(t, c.expected, actual)

		})
	}
}

func TestRoleFrom(t *testing.T) {
	cases := []struct {
		name     string
		literal  string
		expected preauth.Role
	}{
		{name: "case of ADMIN", literal: "ADMIN", expected: preauth.ADMIN},
		{name: "case of USER", literal: "USER", expected: preauth.USER},
		{name: "case of invalid", literal: "Invalid role", expected: preauth.UNDEFINED},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := preauth.RoleFrom(c.literal)
			assert.Equal(t, c.expected, actual)
		})
	}
}
