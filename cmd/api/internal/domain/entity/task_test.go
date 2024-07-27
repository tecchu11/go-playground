package entity_test

import (
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/pkg/errorx"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
