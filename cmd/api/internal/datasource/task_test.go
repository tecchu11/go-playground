package datasource_test

import (
	"context"
	"go-playground/cmd/api/internal/datasource"
	"go-playground/cmd/api/internal/datasource/maindb"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/pkg/errorx"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskAdaptorListTasks(t *testing.T) {
	tests := map[string]struct {
		token    string
		limit    int32
		expected entity.Page[entity.Task]
	}{
		"next is empty": {
			token: "",
			limit: 2,
			expected: entity.Page[entity.Task]{
				Items: []entity.Task{
					{
						ID:        "0191039a-d472-7e9f-9138-7b5e1c400553",
						Content:   "this is test 5",
						CreatedAt: time.Date(2024, 7, 30, 21, 26, 04, 0, time.UTC),
						UpdatedAt: time.Date(2024, 7, 30, 21, 26, 04, 0, time.UTC),
					},
					{
						ID:        "0191039a-cef4-7c15-9b84-525f37ec3f8b",
						Content:   "this is test 4",
						CreatedAt: time.Date(2024, 7, 30, 21, 26, 02, 0, time.UTC),
						UpdatedAt: time.Date(2024, 7, 30, 21, 26, 02, 0, time.UTC),
					},
				},
				HasNext:   true,
				NextToken: "eyJpZCI6IjAxOTEwMmNhLWI1OGItN2I0Ni04ZTI3LWQ2MzQ4NWE3MDU3NCJ9",
			},
		},
		"no next": {
			token: "0190fe59-6618-7811-8b28-a3e67969a4ef",
			limit: 2,
			expected: entity.Page[entity.Task]{
				Items: []entity.Task{
					{
						ID:        "0190fe59-6618-7811-8b28-a3e67969a4ef",
						Content:   "this is test 1",
						CreatedAt: time.Date(2024, 7, 29, 20, 56, 30, 0, time.UTC),
						UpdatedAt: time.Date(2024, 7, 29, 20, 56, 30, 0, time.UTC),
					},
				},
			},
		},
		"no result": {
			token: "00000000-0000-1000-8000-000000000000",
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
				CreatedAt: time.Date(2024, 7, 29, 20, 56, 30, 0, time.UTC),
				UpdatedAt: time.Date(2024, 7, 29, 20, 56, 30, 0, time.UTC),
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
