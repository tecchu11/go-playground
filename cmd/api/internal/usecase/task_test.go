package usecase_test

import (
	"context"
	"errors"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/usecase"
	"go-playground/pkg/errorx"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTaskUseCase_ListTasks(t *testing.T) {
	mockTaskRepo := MockTaskRepository{}
	useCase := usecase.NewTaskUseCase(&mockTaskRepo, nil)
	expected := entity.Page[entity.Task]{
		Items: []entity.Task{
			{ID: "task-id-1"}, {ID: "task-id-2"},
		},
		HasNext:   true,
		NextToken: "task-id-3",
	}
	mockTaskRepo.On("ListTasks", context.Background(), "task-id-1", int32(2)).Return(expected, nil)

	actuaPage, acutalErr := useCase.ListTasks(context.Background(), "eyJpZCI6InRhc2staWQtMSJ9Cg==", 2)

	assert.NoError(t, acutalErr)
	assert.Equal(t, expected, actuaPage)
}

func TestTaskUseCase_ListTasks_LimitIsZero(t *testing.T) {
	mockTaskRepo := MockTaskRepository{}
	useCase := usecase.NewTaskUseCase(&mockTaskRepo, nil)
	expected := entity.Page[entity.Task]{
		Items: []entity.Task{
			{ID: "task-id-1"}, {ID: "task-id-2"},
		},
	}
	mockTaskRepo.On("ListTasks", context.Background(), "task-id-1", int32(10)).Return(expected, nil)

	actuaPage, acutalErr := useCase.ListTasks(context.Background(), "eyJpZCI6InRhc2staWQtMSJ9Cg==", 0)

	assert.NoError(t, acutalErr)
	assert.Equal(t, expected, actuaPage)
}

func TestTaskUseCase_ListTasks_CursorIsInValid(t *testing.T) {
	useCase := usecase.NewTaskUseCase(nil, nil)
	actuaPage, acutalErr := useCase.ListTasks(context.Background(), "invalid", 0)

	var err *errorx.Error
	require.ErrorAs(t, acutalErr, &err)
	assert.Equal(t, "failed to decode task cursor token", err.Msg())
	assert.Equal(t, 400, err.HTTPStatus())
	assert.Equal(t, slog.LevelWarn, err.Level())
	assert.Zero(t, actuaPage)
}

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

	id, err := useCase.CreateTask(context.Background(), "do test")

	taskRepo.AssertExpectations(t)
	assert.NoError(t, err)
	assert.NotZero(t, id)
}

func TestTask_CreateTask_ErrorWhenContentIsBlank(t *testing.T) {
	useCase := usecase.NewTaskUseCase(nil, nil)

	id, err := useCase.CreateTask(context.Background(), " ")

	var myErr *errorx.Error
	assert.ErrorAs(t, err, &myErr)
	assert.Zero(t, id)
}

func TestTask_CreateTask_ErrorWhenCreating(t *testing.T) {
	taskRepo := MockTaskRepository{}
	useCase := usecase.NewTaskUseCase(&taskRepo, nil)
	taskRepo.On("Create", context.Background(), mock.AnythingOfType("Task")).Return(errors.New("failed to create"))

	id, err := useCase.CreateTask(context.Background(), "do test")

	assert.EqualError(t, err, "failed to create")
	assert.Zero(t, id)
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
