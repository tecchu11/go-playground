package testhelper

import (
	"testing"

	"github.com/google/uuid"
)

func UUIDFromString(t *testing.T, value string) uuid.UUID {
	id, err := uuid.Parse(value)
	if err != nil {
		t.Fatalf("uuid from string: %v", err)
	}
	return id
}
