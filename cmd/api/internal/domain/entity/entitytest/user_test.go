package entitytest_test

import (
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/domain/entity/entitytest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestTestUser(t *testing.T) {
	tests := map[string]struct {
		overrides []func(*entity.User)
		want      entity.User
	}{
		"success with default": {
			want: entity.User{
				Sub:           "test-sub",
				GivenName:     "Swaniawski",
				FamilyName:    "Sons",
				Email:         "Emmett.Veum61@example.com",
				EmailVerified: true,
			},
		},
		"success with override": {
			overrides: []func(*entity.User){
				func(u *entity.User) {
					u.Sub = "override-sub"
					u.GivenName = "override-givenName"
					u.FamilyName = "override-familyName"
					u.Email = "Miguel41@example.com"
					u.EmailVerified = false
				},
			},
			want: entity.User{
				Sub:           "override-sub",
				GivenName:     "override-givenName",
				FamilyName:    "override-familyName",
				Email:         "Miguel41@example.com",
				EmailVerified: false,
			},
		},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			got := entitytest.TestUser(t, v.overrides...)

			diff := cmp.Diff(v.want, got, cmpopts.IgnoreFields(got, "ID", "CreatedAt", "UpdatedAt"))
			assert.Empty(t, diff)
		})
	}
}
