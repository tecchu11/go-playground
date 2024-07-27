package repository

import (
	"context"
	"go-playground/cmd/api/internal/domain/entity"
)

// TaskRepository is interface to interact task datasource.
type TaskRepository interface {
	// FindByID find task by given id. Error will be returned if task is not found.
	FindByID(context.Context, entity.TaskID) (entity.Task, error)
	//
	Create(context.Context, entity.Task) error
	Update(context.Context, entity.Task) error
}
