package apperr

import (
	"database/sql"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithLevel(t *testing.T) {
	err := New("text", "msg", WithLevel(slog.LevelInfo))
	appErr := MustAppErr(t, err)

	assert.Equal(t, slog.LevelInfo, *appErr.logLevel)
}

func TestWithCause(t *testing.T) {
	err := New("text", "msg", WithCause(sql.ErrNoRows))
	appErr := MustAppErr(t, err)

	assert.Equal(t, sql.ErrNoRows, appErr.cause)
}
