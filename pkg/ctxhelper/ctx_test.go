package ctxhelper_test

import (
	"context"
	"go-playground/pkg/ctxhelper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubject(t *testing.T) {
	type want struct {
		sub string
		ok  bool
	}
	tests := map[string]struct {
		setup func(context.Context) context.Context
		want  want
	}{
		"success: get subject from context after authorization": {
			setup: func(ctx context.Context) context.Context {
				return ctxhelper.WithSubject(ctx, "90a410d4-025c-4639-8b46-132277ae13e7")
			},
			want: want{
				sub: "90a410d4-025c-4639-8b46-132277ae13e7",
				ok:  true,
			},
		},
		"failure: missing subject from context before authorization": {
			setup: func(ctx context.Context) context.Context {
				return ctx
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := tc.setup(t.Context())

			got, ok := ctxhelper.Subject(ctx)

			assert.Equal(t, got, tc.want.sub)
			assert.Equal(t, ok, tc.want.ok)
		})
	}
}
