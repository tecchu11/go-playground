package httpx_test

import (
	"go-playground/pkg/httpx"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQueryInt32(t *testing.T) {
	tests := map[string]struct {
		r           *http.Request
		expected    int32
		expectError bool
	}{
		"success": {
			r:        httptest.NewRequest("", "/index?key=100", nil),
			expected: 100,
		},
		"parse error": {
			r:           httptest.NewRequest("", "/index?key=value", nil),
			expectError: true,
		},
		"missing": {
			r: httptest.NewRequest("", "/index", nil),
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			actual, err := httpx.QueryInt32(v.r, "key")

			if actual != v.expected {
				t.Fatalf("Actual is %d and expected %d", actual, v.expected)
			}
			if v.expectError && err == nil {
				t.Fatal("err must not be nil")
			}
			if !v.expectError && err != nil {
				t.Fatal("err must be nil")
			}
		})
	}
}
