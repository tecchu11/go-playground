package collection_test

import (
	"go-playground/pkg/collection"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSMap(t *testing.T) {
	tests := map[string]struct {
		in   []int
		want []string
	}{
		"s is nil": {
			want: []string{},
		},
		"s is empty": {
			in:   []int{},
			want: []string{},
		},
		"s is [1, 2, 3]": {
			in:   []int{1, 2, 3},
			want: []string{"1", "2", "3"},
		},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			got := collection.SMap(v.in, func(elem int) string {
				return strconv.Itoa(elem)
			})
			assert.Equal(t, v.want, got)
		})
	}
}
