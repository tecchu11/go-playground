package apperr

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStackTraceNRAttribute(t *testing.T) {
	st := caller(1, 1)

	attr := st.NRAttribute()
	got := attr["stacktrace"]

	assert.Regexp(t, regexp.MustCompile(`^go-playground/pkg/apperr.TestStackTraceNRAttribute\(.+/pkg/apperr/stacktrace_test.go:\d+\)\n$`), got)
}
