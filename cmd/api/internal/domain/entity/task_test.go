package entity_test

import (
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/pkg/errorx"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTask(t *testing.T) {
	actual, actualErr := entity.NewTask("do test")
	assert.NoError(t, actualErr)
	assert.NotZero(t, actual.ID)
	assert.Equal(t, "do test", actual.Content)
	assert.True(t, actual.CreatedAt.Equal(actual.UpdatedAt))
}

func TestNewTask_ErrorWhenContentIsBlank(t *testing.T) {
	actual, err := entity.NewTask("  ")

	assert.Zero(t, actual)
	var myErr *errorx.Error
	assert.ErrorAs(t, err, &myErr)
	assert.Equal(t, 400, myErr.HTTPStatus())
	assert.Equal(t, slog.LevelInfo, myErr.Level())
}

func TestTask_UpdateContent(t *testing.T) {
	task := entity.Task{
		ID:        "test-id",
		Content:   "do test",
		CreatedAt: time.Date(2024, 7, 25, 0, 0, 0, 0, time.Local),
		UpdatedAt: time.Date(2024, 7, 25, 0, 0, 0, 0, time.Local),
	}
	expected := entity.Task{
		ID:        "test-id",
		Content:   "done test",
		CreatedAt: time.Date(2024, 7, 25, 0, 0, 0, 0, time.Local),
		UpdatedAt: time.Date(2024, 7, 25, 0, 0, 0, 0, time.Local),
	}
	err := task.UpdateContent("done test")
	assert.NoError(t, err)
	assert.Equal(t, expected, task)
}

func TestTask_UpdateContent_ErrorWhenContentIsBlank(t *testing.T) {
	task := entity.Task{
		ID:        "test-id",
		Content:   "do test",
		CreatedAt: time.Date(2024, 7, 25, 0, 0, 0, 0, time.Local),
		UpdatedAt: time.Date(2024, 7, 25, 0, 0, 0, 0, time.Local),
	}
	err := task.UpdateContent(" ")
	var myErr *errorx.Error
	assert.ErrorAs(t, err, &myErr)
	assert.Equal(t, 400, myErr.HTTPStatus())
	assert.Equal(t, slog.LevelInfo, myErr.Level())
}

func TestDecodeTaskCursor(t *testing.T) {
	tests := map[string]struct {
		token    string
		expected entity.TaskCursor
	}{
		"token is empty": {},
		"token holds 71113d46-53f1-4ab7-a1c7-0074e707b764": {
			token:    "eyJpZCI6IjcxMTEzZDQ2LTUzZjEtNGFiNy1hMWM3LTAwNzRlNzA3Yjc2NCJ9Cg==",
			expected: entity.TaskCursor{ID: "71113d46-53f1-4ab7-a1c7-0074e707b764"},
		},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			actual, actualErr := entity.DecodeTaskCursor(v.token)
			require.NoError(t, actualErr)
			assert.Equal(t, v.expected, actual)
		})
	}
}

func TestDecodeTaskCursor_Error(t *testing.T) {
	tests := map[string]struct {
		token    string
		expected struct {
			msg    string
			level  slog.Level
			status int
		}
	}{
		"failed to decode base64": {
			token: "not base64",
			expected: struct {
				msg    string
				level  slog.Level
				status int
			}{
				msg:    "failed to decode task cursor token",
				level:  slog.LevelWarn,
				status: 400,
			},
		},
		"failed to unmarshal json": {
			token: "e10K",
			expected: struct {
				msg    string
				level  slog.Level
				status int
			}{
				msg:    "failed to decode task cursor token",
				level:  slog.LevelWarn,
				status: 400,
			},
		},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			actual, actualErr := entity.DecodeTaskCursor(v.token)

			var err *errorx.Error
			require.ErrorAs(t, actualErr, &err)
			assert.Equal(t, v.expected.msg, err.Msg())
			assert.Equal(t, v.expected.level, err.Level())
			assert.Equal(t, v.expected.status, err.HTTPStatus())
			assert.Zero(t, actual)
		})
	}
}
