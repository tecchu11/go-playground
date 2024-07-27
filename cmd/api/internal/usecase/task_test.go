package usecase_test

import (
	"context"
	"errors"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/usecase"
	"go-playground/pkg/errorx"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTaskUseCase_FindTaskByID(t *testing.T) {
	mockTaskRepo := MockTaskRepository{}
	useCase := usecase.NewTaskUseCase(&mockTaskRepo, nil)
	mockTaskRepo.On("FindByID", context.Background(), "task-id").Return(entity.Task{ID: "task-id"}, nil)

	actualTask, actualErr := useCase.FindTaskByID(context.Background(), "task-id")

	mockTaskRepo.AssertExpectations(t)
	assert.NoError(t, actualErr)
	assert.Equal(t, entity.Task{ID: "task-id"}, actualTask)
}

func TestTaskUseCase_FindTaskByID_ErrorWhenTaskIsNotFound(t *testing.T) {
	taskRepo := MockTaskRepository{}
	useCase := usecase.NewTaskUseCase(&taskRepo, nil)
	taskRepo.On("FindByID", context.Background(), "missing-id").Return(entity.Task{}, errors.New("missing"))

	actualTask, actualErr := useCase.FindTaskByID(context.Background(), "missing-id")

	taskRepo.AssertExpectations(t)
	assert.Equal(t, errors.New("missing"), actualErr)
	assert.Zero(t, actualTask)
}

func TestTask_CreateTask(t *testing.T) {
	taskRepo := MockTaskRepository{}
	useCase := usecase.NewTaskUseCase(&taskRepo, nil)
	taskRepo.On("Create", context.Background(), mock.AnythingOfType("Task")).Return(nil)

	err := useCase.CreateTask(context.Background(), "do test")

	taskRepo.AssertExpectations(t)
	assert.NoError(t, err)
}

func TestTask_CreateTask_ErrorWhenContentIsBlank(t *testing.T) {
	useCase := usecase.NewTaskUseCase(nil, nil)

	err := useCase.CreateTask(context.Background(), " ")

	var myErr *errorx.Error
	assert.ErrorAs(t, err, &myErr)
}

func TestTask_CreateTask_ErrorWhenCreating(t *testing.T) {
	taskRepo := MockTaskRepository{}
	useCase := usecase.NewTaskUseCase(&taskRepo, nil)
	taskRepo.On("Create", context.Background(), mock.AnythingOfType("Task")).Return(errors.New("failed to create"))

	err := useCase.CreateTask(context.Background(), "do test")

	assert.EqualError(t, err, "failed to create")
}

func TestTask_UpdateTask(t *testing.T) {
	taskRepo := MockTaskRepository{}
	transaction := MockTransactionRepository{}
	useCase := usecase.NewTaskUseCase(&taskRepo, &transaction)
	taskRepo.On("FindByID", context.Background(), "task-id").Return(entity.Task{ID: "task-id", Content: "content"}, nil)
	taskRepo.On("Update", context.Background(), entity.Task{ID: "task-id", Content: "new-content"}).Return(nil)

	actualErr := useCase.UpdateTask(context.Background(), "task-id", "new-content")

	taskRepo.AssertExpectations(t)
	assert.NoError(t, actualErr)
}

func TestTask_UpdateTask_ErrorWhenTaskIsNotFound(t *testing.T) {
	taskRepo := MockTaskRepository{}
	transaction := MockTransactionRepository{}
	useCase := usecase.NewTaskUseCase(&taskRepo, &transaction)
	taskRepo.On("FindByID", context.Background(), "invalid-id").Return(entity.Task{}, errors.New("missing task"))

	actualErr := useCase.UpdateTask(context.Background(), "invalid-id", "new-content")

	taskRepo.AssertExpectations(t)
	assert.EqualError(t, actualErr, "missing task")
}

func TestTask_UpdateTask_ErrorWhenNewContentIsBlank(t *testing.T) {
	taskRepo := MockTaskRepository{}
	transaction := MockTransactionRepository{}
	useCase := usecase.NewTaskUseCase(&taskRepo, &transaction)
	taskRepo.On("FindByID", context.Background(), "task-id").Return(entity.Task{ID: "task-id", Content: "content"}, nil)

	actualErr := useCase.UpdateTask(context.Background(), "task-id", " ")

	taskRepo.AssertExpectations(t)
	var myErr *errorx.Error
	assert.ErrorAs(t, actualErr, &myErr)
}

func TestTask_UpdateTask_ErrorWhenUpdating(t *testing.T) {
	taskRepo := MockTaskRepository{}
	transaction := MockTransactionRepository{}
	useCase := usecase.NewTaskUseCase(&taskRepo, &transaction)
	taskRepo.On("FindByID", context.Background(), "task-id").Return(entity.Task{ID: "task-id", Content: "content"}, nil)
	taskRepo.On("Update", context.Background(), entity.Task{ID: "task-id", Content: "new-content"}).Return(errors.New("failed to update"))

	actualErr := useCase.UpdateTask(context.Background(), "task-id", "new-content")

	assert.EqualError(t, actualErr, "failed to update")
}
