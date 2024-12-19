package usecase_test

import (
	"context"
	"database/sql"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/usecase"
	"go-playground/pkg/apperr"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTaskUseCase_ListTasks(t *testing.T) {
	type input struct {
		ctx   context.Context
		next  string
		limit int32
	}
	type want struct {
		tasks   entity.Page[entity.Task]
		err     string
		errCode apperr.Code
	}
	type setup func(t *testing.T) *usecase.TaskUseCase
	tests := map[string]struct {
		input input
		setup setup
		want  want
	}{
		"success with param": {
			input: input{ctx: context.Background(), next: "ewogICJpZCI6ICIwMTkzZGQ5Ni00N2FhLTc1NWItODA2ZS0wYjIyZDZmMTg0OWIiCn0K", limit: 2},
			setup: func(t *testing.T) *usecase.TaskUseCase {
				mck := new(MockTaskRepository)
				mck.
					On("ListTasks", context.Background(), "0193dd96-47aa-755b-806e-0b22d6f1849b", int32(2)).
					Return(entity.Page[entity.Task]{
						Items:     []entity.Task{{ID: "0193dd97-123b-7bbe-8229-fa6c91b07a0e"}, {ID: "0193dd97-2b48-711a-b67a-8e9dd44a2dbb"}},
						HasNext:   true,
						NextToken: "0193dd97-565f-755f-8161-e3265eb7a5df",
					}, nil)
				u := usecase.NewTaskUseCase(mck, &MockTransactionRepository{})
				return u
			},
			want: want{
				tasks: entity.Page[entity.Task]{
					Items:     []entity.Task{{ID: "0193dd97-123b-7bbe-8229-fa6c91b07a0e"}, {ID: "0193dd97-2b48-711a-b67a-8e9dd44a2dbb"}},
					HasNext:   true,
					NextToken: "0193dd97-565f-755f-8161-e3265eb7a5df",
				},
			},
		},
		"success without next": {
			input: input{ctx: context.Background(), next: "", limit: 2},
			setup: func(t *testing.T) *usecase.TaskUseCase {
				mck := new(MockTaskRepository)
				mck.
					On("ListTasks", context.Background(), "", int32(2)).
					Return(entity.Page[entity.Task]{
						Items:     []entity.Task{{ID: "0193dd97-123b-7bbe-8229-fa6c91b07a0e"}, {ID: "0193dd97-2b48-711a-b67a-8e9dd44a2dbb"}},
						HasNext:   true,
						NextToken: "0193dd97-565f-755f-8161-e3265eb7a5df",
					}, nil)
				u := usecase.NewTaskUseCase(mck, &MockTransactionRepository{})
				return u
			},
			want: want{
				tasks: entity.Page[entity.Task]{
					Items:     []entity.Task{{ID: "0193dd97-123b-7bbe-8229-fa6c91b07a0e"}, {ID: "0193dd97-2b48-711a-b67a-8e9dd44a2dbb"}},
					HasNext:   true,
					NextToken: "0193dd97-565f-755f-8161-e3265eb7a5df",
				},
			},
		},
		"success without limit": {
			input: input{ctx: context.Background(), next: "ewogICJpZCI6ICIwMTkzZGRhMC1iNWU1LTc2NjctYWI2OS05YjU1YjRiN2UyMGIiCn0K"},
			setup: func(t *testing.T) *usecase.TaskUseCase {
				mck := new(MockTaskRepository)
				mck.
					On("ListTasks", context.Background(), "0193dda0-b5e5-7667-ab69-9b55b4b7e20b", int32(10)).
					Return(entity.Page[entity.Task]{
						Items:     []entity.Task{{ID: "0193dda2-ce96-7333-87a5-94ca410d63e4"}, {ID: "0193dd97-2b48-711a-b67a-8e9dd44a2dbb"}},
						HasNext:   false,
						NextToken: "",
					}, nil)
				return usecase.NewTaskUseCase(mck, &MockTransactionRepository{})
			},
			want: want{
				tasks: entity.Page[entity.Task]{
					Items:     []entity.Task{{ID: "0193dda2-ce96-7333-87a5-94ca410d63e4"}, {ID: "0193dd97-2b48-711a-b67a-8e9dd44a2dbb"}},
					HasNext:   false,
					NextToken: "",
				},
			},
		},
		"failure invalid token": {
			input: input{ctx: context.Background(), next: "invalid"},
			setup: func(t *testing.T) *usecase.TaskUseCase {
				return usecase.NewTaskUseCase(nil, nil)
			},
			want: want{err: "decode task cursor by base64: illegal base64 data at input byte 4", errCode: apperr.CodeInvalidArgument},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			u := tc.setup(t)

			got, err := u.ListTasks(tc.input.ctx, tc.input.next, tc.input.limit)

			if tc.want.err != "" {
				assert.Zero(t, got)
				assert.EqualError(t, err, tc.want.err)
				assert.True(t, apperr.IsCode(err, tc.want.errCode))
			} else {
				assert.Equal(t, tc.want.tasks, got)
				assert.NoError(t, err)
			}
		})
	}
}

func TestTaskUseCase_FindTaskByID(t *testing.T) {
	type input struct {
		ctx context.Context
		id  string
	}
	type want struct {
		task    entity.Task
		err     string
		errCode apperr.Code
	}
	type setup func(t *testing.T) *usecase.TaskUseCase
	tests := map[string]struct {
		input input
		setup setup
		want  want
	}{
		"success": {
			input: input{ctx: context.Background(), id: "0193ddaa-6fdb-7bb6-b6ca-3ee5f131f1f4"},
			setup: func(t *testing.T) *usecase.TaskUseCase {
				mck := new(MockTaskRepository)
				mck.
					On("FindByID", context.Background(), "0193ddaa-6fdb-7bb6-b6ca-3ee5f131f1f4").
					Return(entity.Task{ID: "0193ddaa-6fdb-7bb6-b6ca-3ee5f131f1f4"}, nil)
				return usecase.NewTaskUseCase(mck, nil)
			},
			want: want{task: entity.Task{ID: "0193ddaa-6fdb-7bb6-b6ca-3ee5f131f1f4"}},
		},
		"failure not found task": {
			input: input{ctx: context.Background(), id: "0193ddb0-0054-777d-a60b-cee300725c64"},
			setup: func(t *testing.T) *usecase.TaskUseCase {
				mck := new(MockTaskRepository)
				mck.
					On("FindByID", context.Background(), "0193ddb0-0054-777d-a60b-cee300725c64").
					Return(entity.Task{}, apperr.New("find task", "task not found", apperr.WithCause(sql.ErrNoRows), apperr.CodeNotFound))
				return usecase.NewTaskUseCase(mck, nil)
			},
			want: want{err: "find task: sql: no rows in result set", errCode: apperr.CodeNotFound},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			u := tc.setup(t)

			got, err := u.FindTaskByID(tc.input.ctx, tc.input.id)

			if tc.want.err != "" {
				assert.Zero(t, got)
				assert.EqualError(t, err, tc.want.err)
				assert.True(t, apperr.IsCode(err, tc.want.errCode))
			} else {
				assert.Equal(t, tc.want.task, got)
				assert.NoError(t, err)
			}
		})
	}
}

func TestTask_CreateTask(t *testing.T) {
	type input struct {
		ctx     context.Context
		content string
	}
	type setup func(*testing.T, input) *usecase.TaskUseCase
	type want struct {
		err     string
		errCode apperr.Code
	}
	tests := map[string]struct {
		input input
		setup setup
		want  want
	}{
		"success": {
			input: input{ctx: context.Background(), content: "do test"},
			setup: func(t *testing.T, i input) *usecase.TaskUseCase {
				mck := new(MockTaskRepository)
				matcher := mock.MatchedBy(func(task entity.Task) bool {
					diff := cmp.Diff(task, entity.Task{Content: i.content}, cmpopts.IgnoreFields(entity.Task{}, "ID", "CreatedAt", "UpdatedAt"))
					require.Empty(t, diff)
					require.NotZero(t, task.ID)
					require.NotZero(t, task.CreatedAt)
					require.NotZero(t, task.UpdatedAt)
					return true
				})
				mck.
					On("Create", context.Background(), matcher).
					Return(nil)
				return usecase.NewTaskUseCase(mck, nil)
			},
			want: want{},
		},
		"failure to create task when repository returned error": {
			input: input{ctx: context.Background(), content: "do test"},
			setup: func(t *testing.T, i input) *usecase.TaskUseCase {
				mck := new(MockTaskRepository)
				matcher := mock.MatchedBy(func(task entity.Task) bool {
					diff := cmp.Diff(task, entity.Task{Content: i.content}, cmpopts.IgnoreFields(entity.Task{}, "ID", "CreatedAt", "UpdatedAt"))
					require.Empty(t, diff)
					require.NotZero(t, task.ID)
					require.NotZero(t, task.CreatedAt)
					require.NotZero(t, task.UpdatedAt)
					return true
				})
				mck.On("Create", context.Background(), matcher).Return(apperr.New("failed to save", "failed to create task"))
				return usecase.NewTaskUseCase(mck, nil)
			},
			want: want{err: "failed to save", errCode: apperr.CodeInternal},
		},
		"failure to create task when content is blank": {
			input: input{ctx: context.Background()},
			setup: func(t *testing.T, i input) *usecase.TaskUseCase { return usecase.NewTaskUseCase(nil, nil) },
			want:  want{err: "task content must be non empty", errCode: apperr.CodeInvalidArgument},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			u := tc.setup(t, tc.input)

			got, err := u.CreateTask(tc.input.ctx, tc.input.content)

			if tc.want.err != "" {
				assert.Zero(t, got)
				assert.EqualError(t, err, tc.want.err)
				assert.True(t, apperr.IsCode(err, tc.want.errCode))
			} else {
				assert.NotZero(t, got)
				assert.NoError(t, err)
			}
		})
	}
}

func TestTask_UpdateTask(t *testing.T) {
	type input struct {
		ctx         context.Context
		id, content string
	}
	type setup func(*testing.T) *usecase.TaskUseCase
	type want struct {
		err     string
		errCode apperr.Code
	}
	tests := map[string]struct {
		input input
		setup setup
		want  want
	}{
		"success to update task": {
			input: input{ctx: context.Background(), id: "0193df27-fa0e-7889-9563-2c265d14d185", content: "done test"},
			setup: func(t *testing.T) *usecase.TaskUseCase {
				mck := new(MockTaskRepository)
				mck.On("FindByID", context.Background(), "0193df27-fa0e-7889-9563-2c265d14d185").Return(entity.Task{
					ID:        "0193df27-fa0e-7889-9563-2c265d14d185",
					Content:   "do test",
					CreatedAt: time.Date(2024, 12, 19, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, 12, 19, 0, 0, 0, 0, time.UTC),
				}, nil)
				matcher := mock.MatchedBy(func(task entity.Task) bool {
					diff := cmp.Diff(task, entity.Task{
						ID:        "0193df27-fa0e-7889-9563-2c265d14d185",
						Content:   "done test",
						CreatedAt: time.Date(2024, 12, 19, 0, 0, 0, 0, time.UTC),
					}, cmpopts.IgnoreFields(entity.Task{}, "UpdatedAt"))
					require.Empty(t, diff)
					require.Greater(t, task.UpdatedAt, time.Date(2024, 12, 19, 0, 0, 0, 0, time.UTC))
					return true
				})
				mck.On("Update", context.Background(), matcher).Return(nil)
				return usecase.NewTaskUseCase(mck, &MockTransactionRepository{})
			},
			want: want{},
		},
		"failure to update task when repository returned error on updating": {
			input: input{ctx: context.Background(), id: "0193df28-348c-777a-b989-0009a50791e7", content: "done test"},
			setup: func(t *testing.T) *usecase.TaskUseCase {
				mck := new(MockTaskRepository)
				mck.On("FindByID", context.Background(), "0193df28-348c-777a-b989-0009a50791e7").Return(entity.Task{
					ID:        "0193df28-348c-777a-b989-0009a50791e7",
					Content:   "do test",
					CreatedAt: time.Date(2024, 12, 19, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, 12, 19, 0, 0, 0, 0, time.UTC),
				}, nil)
				matcher := mock.MatchedBy(func(task entity.Task) bool {
					diff := cmp.Diff(task, entity.Task{
						ID:        "0193df28-348c-777a-b989-0009a50791e7",
						Content:   "done test",
						CreatedAt: time.Date(2024, 12, 19, 0, 0, 0, 0, time.UTC),
					}, cmpopts.IgnoreFields(entity.Task{}, "UpdatedAt"))
					require.Empty(t, diff)
					require.Greater(t, task.UpdatedAt, time.Date(2024, 12, 19, 0, 0, 0, 0, time.UTC))
					return true
				})
				mck.On("Update", context.Background(), matcher).Return(apperr.New("failed to save", "failed to save"))
				return usecase.NewTaskUseCase(mck, &MockTransactionRepository{})
			},
			want: want{err: "failed to save", errCode: apperr.CodeInternal},
		},
		"failure to update task when content is empty": {
			input: input{ctx: context.Background(), id: "0193df31-158a-7eee-b12e-3bd316ea15dd"},
			setup: func(t *testing.T) *usecase.TaskUseCase {
				mck := new(MockTaskRepository)
				mck.On("FindByID", context.Background(), "0193df31-158a-7eee-b12e-3bd316ea15dd").Return(entity.Task{
					ID:        "0193df31-158a-7eee-b12e-3bd316ea15dd",
					Content:   "do test",
					CreatedAt: time.Date(2024, 12, 19, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, 12, 19, 0, 0, 0, 0, time.UTC),
				}, nil)
				return usecase.NewTaskUseCase(mck, &MockTransactionRepository{})
			},
			want: want{err: "task content must be non empty", errCode: apperr.CodeInvalidArgument},
		},
		"failure to update task when task not found": {
			input: input{ctx: context.Background(), id: "0193df32-f54d-7330-a242-bc72ae85d7b4", content: "not found"},
			setup: func(t *testing.T) *usecase.TaskUseCase {
				mck := new(MockTaskRepository)
				mck.On("FindByID", context.Background(), "0193df32-f54d-7330-a242-bc72ae85d7b4").Return(entity.Task{}, apperr.New("find task", "task not found", apperr.CodeNotFound, apperr.WithCause(sql.ErrNoRows)))
				return usecase.NewTaskUseCase(mck, &MockTransactionRepository{})
			},
			want: want{err: "find task: sql: no rows in result set", errCode: apperr.CodeNotFound},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			u := tc.setup(t)

			err := u.UpdateTask(tc.input.ctx, tc.input.id, tc.input.content)

			if tc.want.err != "" {
				assert.EqualError(t, err, tc.want.err)
				assert.True(t, apperr.IsCode(err, tc.want.errCode))
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
