package problemdetails_test

import (
	"go-playground/pkg/problemdetails"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProblemDetailsNewJSON(t *testing.T) {
	tests := map[string]struct {
		w        *httptest.ResponseRecorder
		r        *http.Request
		builder  problemdetails.Builder[any]
		expected string
	}{
		"json with default": {
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodGet, "https://example.com/test1", nil),
			builder:  problemdetails.New("this is title", 500),
			expected: `{"type":"about:blank","title":"this is title","status":500,"instance":"/test1"}`,
		},
		"json with specified members": {
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodGet, "https://example.com/test2", nil),
			builder:  problemdetails.New("this is title", 500).WithType("https://mydoc.example.com").WithDetail("this is detail"),
			expected: `{"type":"https://mydoc.example.com","title":"this is title","status":500,"detail":"this is detail","instance":"/test2"}`,
		},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			v.builder.JSON(v.w, v.r)
			require.JSONEq(t, v.expected, v.w.Body.String())
			require.Equal(t, 500, v.w.Code)
			require.Equal(t, "application/problem+json", v.w.Header().Get("Content-Type"))
		})
	}
}

func TestProblemDetailsNewWithTypeJSON(t *testing.T) {
	type extension struct {
		Member string `json:"member"`
	}
	tests := map[string]struct {
		w        *httptest.ResponseRecorder
		r        *http.Request
		builder  problemdetails.Builder[*extension]
		expected string
	}{
		"json with default": {
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodGet, "https://example.com/test1", nil),
			builder:  problemdetails.NewWithType[*extension]("this is title", 500),
			expected: `{"type":"about:blank","title":"this is title","status":500,"instance":"/test1"}`,
		},
		"json with specified members": {
			w:        httptest.NewRecorder(),
			r:        httptest.NewRequest(http.MethodGet, "https://example.com/test2", nil),
			builder:  problemdetails.NewWithType[*extension]("this is title", 500).WithType("https://mydoc.example.com").WithDetail("this is detail").WithAdditional(&extension{Member: "this is member"}),
			expected: `{"type":"https://mydoc.example.com","title":"this is title","status":500,"detail":"this is detail","instance":"/test2","additional":{"member":"this is member"}}`,
		},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			v.builder.JSON(v.w, v.r)
			require.JSONEq(t, v.expected, v.w.Body.String())
			require.Equal(t, 500, v.w.Code)
			require.Equal(t, "application/problem+json", v.w.Header().Get("Content-Type"))
		})
	}
}
