package datasource_test

import (
	"context"
	"errors"
	"go-playground/cmd/api/internal/datasource"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBTransactionAdaptorDo(t *testing.T) {
	tests := map[string]struct {
		ctx      context.Context
		action   func(context.Context) error
		expected error
	}{
		"success": {
			ctx:    context.Background(),
			action: func(ctx context.Context) error { return nil },
		},
		"failure action": {
			ctx:      context.Background(),
			action:   func(ctx context.Context) error { return errors.New("action error") },
			expected: errors.New("action error"),
		},
	}

	dbTransactionAdaptor := datasource.NewDBTransactionAdaptor(db)

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			actual := dbTransactionAdaptor.Do(v.ctx, v.action)
			assert.Equal(t, v.expected, actual)
		})
	}
}
