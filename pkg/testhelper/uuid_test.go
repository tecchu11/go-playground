package testhelper_test

import (
	"go-playground/pkg/testhelper"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUUIDFromString(t *testing.T) {
	input := "01931980-804f-7aa5-9f18-b174ad615940"
	want, err := uuid.Parse(input)
	require.NoError(t, err)

	got := testhelper.UUIDFromString(t, input)

	assert.Equal(t, want, got)
}
