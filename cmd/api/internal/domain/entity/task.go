package entity

import (
	"go-playground/pkg/errorx"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TaskID is identifier of task entity.
type TaskID = string

// Task is domain entity.
type Task struct {
	ID        TaskID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewTask creates new task.
// Return error if content is empty.
func NewTask(content string) (Task, error) {
	err := validateTask(content)
	if err != nil {
		return Task{}, err
	}
	id, err := uuid.NewV7()
	if err != nil {
		return Task{}, errorx.NewError("Failed to publish task id", errorx.WithCause(err))
	}
	now := time.Now()
	return Task{
		ID:        id.String(),
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// UpdateContent updates task content.
func (t *Task) UpdateContent(newContent string) error {
	if err := validateTask(newContent); err != nil {
		return err
	}
	t.Content = newContent
	return nil
}

func validateTask(content string) error {
	if trimmed := strings.TrimSpace(content); trimmed == "" {
		return errorx.NewInfo("task content must be non empty", errorx.WithStatus(400))
	}
	return nil
}
