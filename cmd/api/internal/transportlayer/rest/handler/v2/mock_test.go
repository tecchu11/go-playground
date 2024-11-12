package handler_test

import (
	"context"
	"go-playground/cmd/api/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockPinger struct {
	mock.Mock
}

func (mck *MockPinger) PingContext(ctx context.Context) error {
	args := mck.Called(ctx)
	return args.Error(0)
}

type MockTaskInteractor struct {
	mock.Mock
}

func (mck *MockTaskInteractor) ListTasks(ctx context.Context, token string, limit int32) (entity.Page[entity.Task], error) {
	args := mck.Called(ctx, token, limit)
	return args.Get(0).(entity.Page[entity.Task]), args.Error(1)
}

func (mck *MockTaskInteractor) FindTaskByID(ctx context.Context, id string) (entity.Task, error) {
	args := mck.Called(ctx, id)
	return args.Get(0).(entity.Task), args.Error(1)
}

func (mck *MockTaskInteractor) CreateTask(ctx context.Context, content string) (entity.TaskID, error) {
	args := mck.Called(ctx, content)
	return args.Get(0).(string), args.Error(1)
}

func (mck *MockTaskInteractor) UpdateTask(ctx context.Context, id, content string) error {
	args := mck.Called(ctx, id, content)
	return args.Error(0)
}

type MockUserInteractor struct {
	mock.Mock
}

func (mck *MockUserInteractor) CreateUser(ctx context.Context, sub string, givenName, familyName string, email string, emailVerified bool) (uuid.UUID, error) {
	args := mck.Called(ctx, sub, givenName, familyName, email, emailVerified)
	return args.Get(0).(uuid.UUID), args.Error(1)
}
