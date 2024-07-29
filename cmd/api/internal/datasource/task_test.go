package datasource_test

import (
	"context"
	"go-playground/cmd/api/internal/datasource"
	"go-playground/cmd/api/internal/datasource/maindb"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/pkg/timex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskAdaptorFindByID(t *testing.T) {
	tests := map[string]struct {
		id             string
		expectedEntity entity.Task
		expectedErr    error
	}{
		"success": {
			id: "0190f34a-e069-7873-8fe1-fdf871eb3918",
			expectedEntity: entity.Task{
				ID:        "0190f34a-e069-7873-8fe1-fdf871eb3918",
				Content:   "this is test 1",
				CreatedAt: time.Date(2024, 7, 28, 0, 0, 0, 0, timex.JST()),
				UpdatedAt: time.Date(2024, 7, 28, 0, 0, 0, 0, timex.JST()),
			},
			expectedErr: nil,
		},
	}
	adaptor := datasource.NewTaskAdaptor(maindb.New(db))
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			runInTx(t, func(ctx context.Context) {
				actualEntity, actualErr := adaptor.FindByID(ctx, v.id)
				assert.Equal(t, v.expectedEntity, actualEntity)
				assert.Equal(t, v.expectedErr, actualErr)
			})
		})
	}
}

func TestTaskAdaptorCreate(t *testing.T) {
	task := entity.Task{
		ID:      "0190f34a-e069-7873-8fe1-fdf871eb3919",
		Content: "create task",
	}
	adaptor := datasource.NewTaskAdaptor(maindb.New(db))
	runInTx(t, func(ctx context.Context) {
		err := adaptor.Create(ctx, task)
		assert.NoError(t, err)

		_, err = adaptor.FindByID(ctx, task.ID)
		assert.NoError(t, err)
	})
}

func TestTaskAdaptorUpdate(t *testing.T) {
	task := entity.Task{
		ID:      "0190f34a-e069-7873-8fe1-fdf871eb3918",
		Content: "update task",
	}
	adaptor := datasource.NewTaskAdaptor(maindb.New(db))
	runInTx(t, func(ctx context.Context) {
		err := adaptor.Update(ctx, task)
		assert.NoError(t, err)

		actual, err := adaptor.FindByID(ctx, task.ID)
		assert.NoError(t, err)
		assert.Equal(t, task.ID, actual.ID)
		assert.Equal(t, task.Content, actual.Content)
	})
}
