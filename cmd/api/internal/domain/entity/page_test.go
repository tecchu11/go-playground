package entity_test

import (
	"go-playground/cmd/api/internal/domain/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCursorPage(t *testing.T) {
	tests := map[string]struct {
		items    []entity.Task
		expected struct {
			items     []entity.Task
			hasNext   bool
			nextToken string
		}
	}{
		"has next": {
			items: []entity.Task{{ID: "1"}, {ID: "2"}},
			expected: struct {
				items     []entity.Task
				hasNext   bool
				nextToken string
			}{
				items:     []entity.Task{{ID: "1"}},
				hasNext:   true,
				nextToken: "2",
			},
		},
		"no next": {
			items: []entity.Task{{ID: "3"}},
			expected: struct {
				items     []entity.Task
				hasNext   bool
				nextToken string
			}{
				items: []entity.Task{{ID: "3"}},
			},
		},
		"empty items": {},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			page := entity.NewCursorPage(v.items, 1)

			assert.Equal(t, v.expected.items, page.Items)
			assert.Equal(t, v.expected.hasNext, page.HasNext)
			assert.Equal(t, v.expected.nextToken, page.NextToken)
		})
	}

}
