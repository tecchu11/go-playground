package apperr

import (
	"log/slog"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodeString(t *testing.T) {
	tests := map[string]struct {
		input Code
		want  string
	}{
		"CodeUnknown":         {input: CodeUnknown, want: "unknown"},
		"CodeInvalidArgument": {input: CodeInvalidArgument, want: "invalidArgument"},
		"CodeUnAuthn":         {input: CodeUnAuthn, want: "unauthenticated"},
		"CodeUnAuthz":         {input: CodeUnAuthz, want: "unauthorized"},
		"CodeNotFound":        {input: CodeNotFound, want: "notfound"},
		"CodeInternal":        {input: CodeInternal, want: "internalServerError"},
		"CodeCustom":          {input: Code(uint32(100)), want: "Code(100)"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.String()

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCodeStatus(t *testing.T) {
	tests := map[string]struct {
		input Code
		want  int
	}{
		"CodeUnknown":         {input: CodeUnknown, want: http.StatusInternalServerError},
		"CodeInvalidArgument": {input: CodeInvalidArgument, want: http.StatusBadRequest},
		"CodeUnAuthn":         {input: CodeUnAuthn, want: http.StatusUnauthorized},
		"CodeUnAuthz":         {input: CodeUnAuthz, want: http.StatusForbidden},
		"CodeNotFound":        {input: CodeNotFound, want: http.StatusNotFound},
		"CodeInternal":        {input: CodeInternal, want: http.StatusInternalServerError},
		"CodeCustom":          {input: Code(uint32(100)), want: http.StatusInternalServerError},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.status()

			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCodeLevel(t *testing.T) {
	tests := map[string]struct {
		input Code
		want  slog.Level
	}{
		"CodeUnknown":         {input: CodeUnknown, want: slog.LevelError},
		"CodeInvalidArgument": {input: CodeInvalidArgument, want: slog.LevelWarn},
		"CodeUnAuthn":         {input: CodeUnAuthn, want: slog.LevelWarn},
		"CodeUnAuthz":         {input: CodeUnAuthz, want: slog.LevelWarn},
		"CodeNotFound":        {input: CodeNotFound, want: slog.LevelWarn},
		"CodeInternal":        {input: CodeInternal, want: slog.LevelError},
		"CodeCustom":          {input: Code(uint32(100)), want: slog.LevelError},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.input.level()

			assert.Equal(t, tc.want, got)
		})
	}
}
