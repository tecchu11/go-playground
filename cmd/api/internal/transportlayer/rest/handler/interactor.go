package handler

import (
	"context"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/usecase"
)

// TaskInteractor is interface for usecase.TaskUseCase.
type TaskInteractor interface {
	FindTaskByID(context.Context, string) (entity.Task, error)
	CreateTask(context.Context, string) error
	UpdateTask(context.Context, string, string) error
}

var _ TaskInteractor = (*usecase.TaskUseCase)(nil)
