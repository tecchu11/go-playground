package usecase

import (
	"context"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/domain/repository"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type TaskUseCase struct {
	transaction    repository.TransactionRepository
	taskRepository repository.TaskRepository
}

func NewTaskUseCase(taskRepo repository.TaskRepository, transaction repository.TransactionRepository) *TaskUseCase {
	return &TaskUseCase{transaction: transaction, taskRepository: taskRepo}
}

func (u *TaskUseCase) FindTaskByID(ctx context.Context, id string) (entity.Task, error) {
	defer newrelic.FromContext(ctx).StartSegment("usecase/TaskUseCase/FindByTaskID").End()

	task, err := u.taskRepository.FindByID(ctx, id)
	if err != nil {
		return entity.Task{}, err
	}
	return task, nil
}

func (u *TaskUseCase) CreateTask(ctx context.Context, content string) error {
	defer newrelic.FromContext(ctx).StartSegment("usecase/TaskUseCase/CreateTask").End()

	task, err := entity.NewTask(content)
	if err != nil {
		return err
	}
	err = u.taskRepository.Create(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

func (u *TaskUseCase) UpdateTask(ctx context.Context, id string, content string) error {
	defer newrelic.FromContext(ctx).StartSegment("usecase/TaskUseCase/UpdateTask").End()

	return u.transaction.Do(ctx, func(ctx context.Context) error {
		task, err := u.taskRepository.FindByID(ctx, id)
		if err != nil {
			return err
		}
		err = task.UpdateContent(content)
		if err != nil {
			return err
		}
		err = u.taskRepository.Update(ctx, task)
		if err != nil {
			return err
		}
		return nil
	})
}
