package datasource_test

import (
	"context"
	"go-playground/cmd/api/internal/datasource"
	"go-playground/cmd/api/internal/datasource/database"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/pkg/apperr"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskAdaptor_ListTasks(t *testing.T) {
	type input struct {
		limit int32
		token string
	}
	type want struct {
		page entity.Page[entity.Task]
	}
	tests := map[string]struct {
		input input
		want  want
	}{
		"next is empty": {
			input: input{token: "", limit: 2},
			want: want{page: entity.Page[entity.Task]{
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
			}},
		},
		"no next": {
			input: input{token: "0190fe59-6618-7811-8b28-a3e67969a4ef", limit: 2},
			want: want{page: entity.Page[entity.Task]{
				Items: []entity.Task{
					{
						ID:        "0190fe59-6618-7811-8b28-a3e67969a4ef",
						Content:   "this is test 1",
						CreatedAt: time.Date(2024, 7, 29, 20, 56, 30, 0, time.UTC),
						UpdatedAt: time.Date(2024, 7, 29, 20, 56, 30, 0, time.UTC),
					},
				},
			}},
		},
		"no result": {
			input: input{token: "00000000-0000-1000-8000-000000000000", limit: 1},
			want:  want{page: entity.Page[entity.Task]{Items: []entity.Task{}}},
		},
	}
	adaptor := datasource.NewTaskAdaptor(database.New(db))
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runInTx(t, func(ctx context.Context) {
				got, err := adaptor.ListTasks(ctx, tc.input.token, tc.input.limit)

				require.NoError(t, err)
				assert.Equal(t, tc.want.page, got)
			})
		})
	}
}

func TestTaskAdaptor_FindByID(t *testing.T) {
	type input struct {
		id string
	}
	type want struct {
		task    entity.Task
		err     string
		errCode apperr.Code
	}
	tests := map[string]struct {
		input input
		want  want
	}{
		"success": {
			input: input{id: "0190fe59-6618-7811-8b28-a3e67969a4ef"},
			want: want{task: entity.Task{
				ID:        "0190fe59-6618-7811-8b28-a3e67969a4ef",
				Content:   "this is test 1",
				CreatedAt: time.Date(2024, 7, 29, 20, 56, 30, 0, time.UTC),
				UpdatedAt: time.Date(2024, 7, 29, 20, 56, 30, 0, time.UTC),
			}},
		},
		"missing": {
			input: input{id: "0193dd05-ea21-7ee3-9aa8-257efa35307a"},
			want:  want{err: `find task by id "0193dd05-ea21-7ee3-9aa8-257efa35307a": sql: no rows in result set`, errCode: apperr.CodeNotFound},
		},
	}
	adaptor := datasource.NewTaskAdaptor(database.New(db))
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runInTx(t, func(ctx context.Context) {
				got, err := adaptor.FindByID(ctx, tc.input.id)

				if tc.want.err != "" {
					assert.Zero(t, got)
					assert.EqualError(t, err, tc.want.err)
					assert.True(t, apperr.IsCode(err, tc.want.errCode))
				} else {
					assert.Equal(t, tc.want.task, got)
					assert.NoError(t, err)
				}
			})
		})
	}
}

func TestTaskAdaptor_Create(t *testing.T) {
	task := entity.Task{
		ID:      "0190f34a-e069-7873-8fe1-fdf871eb3919",
		Content: "create task",
	}
	adaptor := datasource.NewTaskAdaptor(database.New(db))
	runInTx(t, func(ctx context.Context) {
		err := adaptor.Create(ctx, task)
		assert.NoError(t, err)

		_, err = adaptor.FindByID(ctx, task.ID)
		assert.NoError(t, err)
	})
}

func TestTaskAdaptor_Update(t *testing.T) {
	task := entity.Task{
		ID:      "0190fe59-6618-7811-8b28-a3e67969a4ef",
		Content: "update task",
	}
	adaptor := datasource.NewTaskAdaptor(database.New(db))
	runInTx(t, func(ctx context.Context) {
		err := adaptor.Update(ctx, task)
		assert.NoError(t, err)

		actual, err := adaptor.FindByID(ctx, task.ID)
		assert.NoError(t, err)
		assert.Equal(t, task.ID, actual.ID)
		assert.Equal(t, task.Content, actual.Content)
	})
}
