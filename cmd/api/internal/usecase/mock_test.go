package usecase_test

import (
	"context"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/domain/repository"

	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (mck *MockTaskRepository) ListTasks(ctx context.Context, token entity.TaskID, limit int32) (entity.Page[entity.Task], error) {
	args := mck.Called(ctx, token, limit)
	return args.Get(0).(entity.Page[entity.Task]), args.Error(1)
}

func (mck *MockTaskRepository) FindByID(ctx context.Context, id entity.TaskID) (entity.Task, error) {
	args := mck.Called(ctx, id)
	return args.Get(0).(entity.Task), args.Error(1)
}

func (mck *MockTaskRepository) Create(ctx context.Context, task entity.Task) error {
	args := mck.Called(ctx, task)
	return args.Error(0)
}

func (mck *MockTaskRepository) Update(ctx context.Context, task entity.Task) error {
	args := mck.Called(ctx, task)
	return args.Error(0)
}

type MockTransactionRepository struct{}

func (mck *MockTransactionRepository) Do(ctx context.Context, action func(context.Context) error) error {
	return action(ctx)
}

type MockUserRepository struct {
	mock.Mock
}

func (mck *MockUserRepository) Create(ctx context.Context, user entity.User) error {
	args := mck.Called(ctx, user)
	return args.Error(0)
}

var (
	_ repository.TransactionRepository = (*MockTransactionRepository)(nil)
	_ repository.TaskRepository        = (*MockTaskRepository)(nil)
	_ repository.UserRepository        = (*MockUserRepository)(nil)
)
