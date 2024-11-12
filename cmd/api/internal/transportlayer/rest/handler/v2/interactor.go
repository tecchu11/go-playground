package handler

import (
	"context"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/usecase"

	"github.com/google/uuid"
)

// TaskInteractor is interface for [usecase.TaskUseCase].
type TaskInteractor interface {
	ListTasks(context.Context, string, int32) (entity.Page[entity.Task], error)
	FindTaskByID(context.Context, string) (entity.Task, error)
	CreateTask(context.Context, string) (entity.TaskID, error)
	UpdateTask(context.Context, string, string) error
}

// UserInteractor is interface for [usecase.UserUseCase]
type UserInteractor interface {
	// CreateUser creates user with given information.
	CreateUser(ctx context.Context, sub string, givenName, familyName string, email string, emailVerified bool) (uuid.UUID, error)
}

var (
	_ TaskInteractor = (*usecase.TaskUseCase)(nil)
	_ UserInteractor = (*usecase.UserUseCase)(nil)
)
