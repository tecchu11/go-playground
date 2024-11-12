package testhelper_test

import (
	"go-playground/pkg/testhelper"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpyContextValu(t *testing.T) {
	spy := new(testhelper.SpyContext)
	spy.On("Value", "key").Return("value")

	got := spy.Value("key")
	require.Equal(t, "value", got)
}
