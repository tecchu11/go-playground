package datasource_test

import (
	"context"
	"database/sql"
	"go-playground/cmd/api/internal/datasource"
	"go-playground/pkg/apperr"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBTransactionAdaptor_Do(t *testing.T) {
	type input struct {
		ctx    context.Context
		action func(context.Context) error
	}
	type want struct {
		err     string
		errCode apperr.Code
	}
	tests := map[string]struct {
		input input
		want  want
	}{
		"success": {
			input: input{
				ctx:    context.Background(),
				action: func(ctx context.Context) error { return nil },
			},
		},
		"failure action": {
			input: input{
				ctx: context.Background(),
				action: func(ctx context.Context) error {
					return apperr.New("find something", "not found something.", apperr.WithCause(sql.ErrNoRows), apperr.CodeNotFound)
				},
			},
			want: want{err: "find something: sql: no rows in result set", errCode: apperr.CodeNotFound},
		},
	}

	dbTransactionAdaptor := datasource.NewDBTransactionAdaptor(db)

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := dbTransactionAdaptor.Do(tc.input.ctx, tc.input.action)

			if tc.want.err != "" {
				assert.EqualError(t, err, tc.want.err)
				assert.True(t, apperr.IsCode(err, tc.want.errCode))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
