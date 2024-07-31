package env_test

import (
	"go-playground/pkg/env"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplyString(t *testing.T) {
	tests := map[string]struct {
		lookup func(string) (string, bool)
	}{
		"lookup is nil":          {},
		"lookup is os.LookupEnv": {lookup: os.LookupEnv},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			t.Setenv("FOO", "BAR")
			var err error
			str := env.ApplyString(&err, "FOO", v.lookup)

			require.Equal(t, "BAR", str)
			require.NoError(t, err)
		})
	}
}

func TestApplyString_Error(t *testing.T) {
	var err error
	str := env.ApplyString(&err, "FOO", nil)
	require.Error(t, err)
	require.Zero(t, str)
}
