package entitytest

import (
	"go-playground/cmd/api/internal/domain/entity"
	"testing"
)

// TestUser is helper for [entity.User].
func TestUser(t *testing.T, overrides ...func(*entity.User)) entity.User {
	user, err := entity.NewUser(
		"test-sub",
		"Swaniawski",
		"Sons",
		"Emmett.Veum61@example.com",
		true,
	)
	if err != nil {
		t.Fatalf("test user: %v", err)
	}
	for _, fn := range overrides {
		fn(&user)
	}
	return user
}
