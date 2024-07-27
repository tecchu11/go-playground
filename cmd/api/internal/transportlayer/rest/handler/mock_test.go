package handler_test

import (
	"context"
	"go-playground/cmd/api/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

type MockTaskInteractor struct {
	mock.Mock
}

func (mck *MockTaskInteractor) FindTaskByID(ctx context.Context, id string) (entity.Task, error) {
	args := mck.Called(ctx, id)
	return args.Get(0).(entity.Task), args.Error(1)
}

func (mck *MockTaskInteractor) CreateTask(ctx context.Context, content string) error {
	args := mck.Called(ctx, content)
	return args.Error(0)
}

func (mck *MockTaskInteractor) UpdateTask(ctx context.Context, id, content string) error {
	args := mck.Called(ctx, id, content)
	return args.Error(0)
}
