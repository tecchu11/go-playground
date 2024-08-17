package env_test

import (
	"go-playground/pkg/env/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplierNew(t *testing.T) {
	got := env.New(nil)

	assert.Nil(t, got.Err())
	assert.NotNil(t, got.Lookup())
}

func TestApplier_String(t *testing.T) {
	type in struct {
		key    string
		setEnv func(t *testing.T)
	}
	type want struct {
		v   string
		err string
	}
	tests := map[string]struct {
		in   in
		want want
	}{
		"success": {
			in: in{
				key: "TEST_ENV",
				setEnv: func(t *testing.T) {
					t.Setenv("TEST_ENV", "this is test")
				},
			},
			want: want{
				v: "this is test",
			},
		},
		"missing env": {
			in: in{
				key:    "TEST_ENV",
				setEnv: func(t *testing.T) { /* noop */ },
			},
			want: want{
				err: "missing env TEST_ENV",
			},
		},
	}
	for name, testCase := range tests {
		t.Run(name, func(t *testing.T) {
			testCase.in.setEnv(t)
			applier := env.New(nil)

			gotV := applier.String(testCase.in.key)

			assert.Equal(t, testCase.want.v, gotV)
			if testCase.want.err == "" {
				assert.NoError(t, applier.Err())
			} else {
				assert.EqualError(t, applier.Err(), testCase.want.err)
			}
		})
	}
}
