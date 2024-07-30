package datasource_test

import (
	"context"
	"go-playground/cmd/api/internal/datasource"
	"go-playground/cmd/api/internal/datasource/maindb"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/pkg/errorx"
	"go-playground/pkg/timex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskAdaptorListTasks(t *testing.T) {
	tests := map[string]struct {
		token    string
		limit    int32
		expected entity.CursorPage[string, entity.Task]
	}{
		"default": {
			token: "",
			limit: 2,
			expected: entity.CursorPage[string, entity.Task]{
				Items: []entity.Task{
					{
						ID:        "0190fe59-6618-7811-8b28-a3e67969a4ef",
						Content:   "this is test 1",
						CreatedAt: time.Date(2024, 7, 29, 20, 56, 30, 0, timex.JST()),
						UpdatedAt: time.Date(2024, 7, 29, 20, 56, 30, 0, timex.JST()),
					},
					{
						ID:        "0190fe5b-1f83-7024-a233-c8a18935f5dc",
						Content:   "this is test 2",
						CreatedAt: time.Date(2024, 7, 29, 20, 58, 23, 0, timex.JST()),
						UpdatedAt: time.Date(2024, 7, 29, 20, 58, 23, 0, timex.JST()),
					},
				},
				HasNext:   true,
				NextToken: "019102ca-b58b-7b46-8e27-d63485a70574",
			},
		},
		"no next": {
			token: "0191039a-d472-7e9f-9138-7b5e1c400553",
			limit: 2,
			expected: entity.CursorPage[string, entity.Task]{
				Items: []entity.Task{
					{
						ID:        "0191039a-d472-7e9f-9138-7b5e1c400553",
						Content:   "this is test 5",
						CreatedAt: time.Date(2024, 7, 30, 21, 26, 4, 0, timex.JST()),
						UpdatedAt: time.Date(2024, 7, 30, 21, 26, 4, 0, timex.JST()),
					},
				},
			},
		},
		"empty": {
			token: "019103cf-dcb8-797f-ba3b-78246e157c1c",
			limit: 1,
		},
	}
	adaptor := datasource.NewTaskAdaptor(maindb.New(db))
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			runInTx(t, func(ctx context.Context) {
				actual, err := adaptor.ListTasks(ctx, v.token, v.limit)

				require.NoError(t, err)
				assert.Equal(t, v.expected, actual)
			})
		})
	}
}

func TestTaskAdaptorFindByID(t *testing.T) {
	tests := map[string]struct {
		id             string
		expectedEntity entity.Task
		expectErr      bool
	}{
		"success": {
			id: "0190fe59-6618-7811-8b28-a3e67969a4ef",
			expectedEntity: entity.Task{
				ID:        "0190fe59-6618-7811-8b28-a3e67969a4ef",
				Content:   "this is test 1",
				CreatedAt: time.Date(2024, 7, 29, 20, 56, 30, 0, timex.JST()),
				UpdatedAt: time.Date(2024, 7, 29, 20, 56, 30, 0, timex.JST()),
			},
		},
		"missing": {
			id:        "invalid",
			expectErr: true,
		},
	}
	adaptor := datasource.NewTaskAdaptor(maindb.New(db))
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			runInTx(t, func(ctx context.Context) {
				actualEntity, actualErr := adaptor.FindByID(ctx, v.id)
				assert.Equal(t, v.expectedEntity, actualEntity)
				if v.expectErr {
					assert.Error(t, actualErr)
					var target *errorx.Error
					assert.ErrorAs(t, actualErr, &target)
				} else {
					assert.NoError(t, actualErr)
				}
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
		ID:      "0190fe59-6618-7811-8b28-a3e67969a4ef",
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
