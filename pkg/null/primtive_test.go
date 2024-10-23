package null_test

import (
	"go-playground/pkg/null"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	input := "input"
	want := &input

	got := null.String(input)

	assert.Equal(t, want, got)
}

func TestInt32(t *testing.T) {
	input := int32(1)
	want := &input

	got := null.Int32(input)

	assert.Equal(t, want, got)
}
