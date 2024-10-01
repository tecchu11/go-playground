package entity

import (
	"encoding/base64"
	"encoding/json"
	"go-playground/pkg/errorx"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TaskID is identifier of task entity.
type TaskID = string

// Task is domain entity.
type Task struct {
	ID        TaskID    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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

// Token is implementation of Item[string] interface.
func (t Task) Token() string {
	return t.ID
}

type TaskCursor struct {
	ID string `json:"id"`
}

// EncodeCursor encodes task cursor token.
func (t Task) EncodeCursor() (string, error) {
	c := TaskCursor{ID: t.ID}
	buf, err := json.Marshal(c)
	if err != nil {
		return "", errorx.NewError("failed to create task cursor token", errorx.WithCause(err))
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

// DecodeTaskCursor decodes token to task cursor.
func DecodeTaskCursor(token string) (TaskCursor, error) {
	if token == "" {
		return TaskCursor{}, nil
	}
	var cursor TaskCursor
	buf, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return cursor, errorx.NewWarn("failed to decode task cursor token", errorx.WithCause(err))
	}
	err = json.Unmarshal(buf, &cursor)
	if err != nil {
		return cursor, errorx.NewWarn("failed to decode task cursor token", errorx.WithCause(err))
	}
	return cursor, nil
}
