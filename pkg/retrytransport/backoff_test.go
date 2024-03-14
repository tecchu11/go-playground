package retrytransport_test

import (
	"go-playground/pkg/retrytransport"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultBackOff(t *testing.T) {
	tests := map[string]struct {
		in   int
		want time.Duration
	}{
		"first attempt":    {in: 0, want: 50 * time.Millisecond},
		"second attempt":   {in: 1, want: 2 * 50 * time.Millisecond},
		"third attempt":    {in: 2, want: 4 * 50 * time.Millisecond},
		"numerous attempt": {in: 200, want: 200 * time.Millisecond},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			got := retrytransport.DefaultBackOff(v.in)
			assert.Equal(t, v.want, got)
		})
	}
}
