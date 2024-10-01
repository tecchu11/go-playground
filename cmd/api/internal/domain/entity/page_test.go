package entity_test

import (
	"errors"
	"go-playground/cmd/api/internal/domain/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPage(t *testing.T) {
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
				nextToken: "eyJpZCI6IjIifQ==",
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
			page, err := entity.NewPage(v.items, 1)

			require.NoError(t, err)
			assert.Equal(t, v.expected.items, page.Items)
			assert.Equal(t, v.expected.hasNext, page.HasNext)
			assert.Equal(t, v.expected.nextToken, page.NextToken)
		})
	}
}

type testInValidEntity struct {
	string
}

func (e testInValidEntity) EncodeCursor() (string, error) {
	return "", errors.New("failed to encode cursor")
}

func TestNewPage_Error(t *testing.T) {
	actual, err := entity.NewPage([]testInValidEntity{{"1"}, {"2"}}, 1)

	assert.EqualError(t, err, "failed to encode cursor")
	assert.Zero(t, actual)
}
