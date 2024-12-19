package entity_test

import (
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/pkg/apperr"
	"go-playground/pkg/testhelper"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestNewTask(t *testing.T) {
	type input struct {
		content string
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
		"success to new": {
			input: input{content: "test"},
			want:  want{task: entity.Task{Content: "test"}},
		},
		"failure on validation": {
			input: input{},
			want:  want{err: "task content must be non empty", errCode: apperr.CodeInvalidArgument},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := entity.NewTask(tc.input.content)

			if tc.want.err != "" {
				assert.Zero(t, got)
				assert.EqualError(t, err, tc.want.err)
				assert.True(t, apperr.IsCode(err, tc.want.errCode))
			} else {
				assert.NoError(t, err)
				diff := cmp.Diff(got, tc.want.task, cmpopts.IgnoreFields(entity.Task{}, "ID", "CreatedAt", "UpdatedAt"))
				assert.Empty(t, diff)
				assert.NotZero(t, got.ID)
				assert.NotZero(t, got.CreatedAt)
				assert.NotZero(t, got.UpdatedAt)
			}
		})
	}
}

func TestTask_UpdateContent(t *testing.T) {
	now := time.Now()
	type input struct {
		task       entity.Task
		newContent string
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
		"success to update content": {
			input: input{
				task:       entity.Task{ID: "0193d862-8182-7997-971d-175160c7f7d6", Content: "do test", CreatedAt: now, UpdatedAt: now},
				newContent: "done test",
			},
			want: want{
				task: entity.Task{ID: "0193d862-8182-7997-971d-175160c7f7d6", Content: "done test", CreatedAt: now},
			},
		},
		"failure to update content": {
			input: input{
				task:       entity.Task{ID: "0193d862-8182-7997-971d-175160c7f7d6", Content: "do test", CreatedAt: now, UpdatedAt: now},
				newContent: "",
			},
			want: want{
				err: "task content must be non empty", errCode: apperr.CodeInvalidArgument,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			task := tc.input.task
			err := task.UpdateContent(tc.input.newContent)

			if tc.want.err != "" {
				assert.EqualError(t, err, tc.want.err)
				assert.True(t, apperr.IsCode(err, tc.want.errCode))
			} else {
				assert.NoError(t, err)
				diff := cmp.Diff(task, tc.want.task, cmpopts.IgnoreFields(entity.Task{}, "UpdatedAt"))
				assert.Empty(t, diff)
				assert.Greater(t, task.UpdatedAt, tc.input.task.UpdatedAt)
			}
		})
	}
}

func TestTask_EncodeCursor(t *testing.T) {
	id := testhelper.UUIDFromString(t, "0193da97-5571-7992-8470-e831708ca57a")
	task := entity.Task{ID: id.String()}

	got, err := task.EncodeCursor()

	assert.NoError(t, err)
	assert.Equal(t, "eyJpZCI6IjAxOTNkYTk3LTU1NzEtNzk5Mi04NDcwLWU4MzE3MDhjYTU3YSJ9", got)
}

func TestDecodeTaskCursor(t *testing.T) {
	type input struct{ token string }
	type want struct {
		cursor  entity.TaskCursor
		err     string
		errCode apperr.Code
	}
	tests := map[string]struct {
		input    input
		want     want
		token    string
		expected entity.TaskCursor
	}{
		"token is empty": {},
		"token holds 71113d46-53f1-4ab7-a1c7-0074e707b764": {
			input: input{token: "eyJpZCI6IjcxMTEzZDQ2LTUzZjEtNGFiNy1hMWM3LTAwNzRlNzA3Yjc2NCJ9Cg=="},
			want:  want{cursor: entity.TaskCursor{ID: "71113d46-53f1-4ab7-a1c7-0074e707b764"}},
		},
		"not base64": {
			input: input{token: "not base64"},
			want:  want{err: "decode task cursor by base64: illegal base64 data at input byte 3", errCode: apperr.CodeInvalidArgument},
		},
		"not json": {
			input: input{token: "e10K"},
			want:  want{err: "decode task cursor by json: invalid character ']' looking for beginning of object key string", errCode: apperr.CodeInvalidArgument},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := entity.DecodeTaskCursor(tc.input.token)

			if tc.want.err != "" {
				assert.Zero(t, got)
				assert.EqualError(t, err, tc.want.err)
				assert.True(t, apperr.IsCode(err, tc.want.errCode))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want.cursor, got)
			}
		})
	}
}
