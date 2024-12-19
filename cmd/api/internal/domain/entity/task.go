package entity

import (
	"encoding/base64"
	"encoding/json"
	"go-playground/pkg/apperr"
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
		return Task{}, apperr.New("uuid new v7 for task id", "Failed to create new task.", apperr.WithCause(err))
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
	t.UpdatedAt = time.Now()
	return nil
}

func validateTask(content string) error {
	if trimmed := strings.TrimSpace(content); trimmed == "" {
		return apperr.New("task content must be non empty", "Task content must be non empty", apperr.CodeInvalidArgument)
	}
	return nil
}

type TaskCursor struct {
	ID string `json:"id"`
}

// EncodeCursor encodes task cursor token.
func (t Task) EncodeCursor() (string, error) {
	c := TaskCursor{ID: t.ID}
	buf, err := json.Marshal(c)
	if err != nil {
		return "", apperr.New("marshal task cursor", "Failed to create task metadata", apperr.WithCause(err))
	}
	return base64.StdEncoding.EncodeToString(buf), nil
}

// DecodeTaskCursor decodes token to task cursor.
func DecodeTaskCursor(token string) (TaskCursor, error) {
	if token == "" {
		return TaskCursor{}, nil
	}
	b, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return TaskCursor{}, apperr.New("decode task cursor by base64", "invalid task cursor", apperr.WithCause(err), apperr.CodeInvalidArgument)
	}
	var cursor TaskCursor
	err = json.Unmarshal(b, &cursor)
	if err != nil {
		return TaskCursor{}, apperr.New("decode task cursor by json", "invalid task cursor", apperr.WithCause(err), apperr.CodeInvalidArgument)
	}
	return cursor, nil
}
