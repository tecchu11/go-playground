package handler_test

import (
	"context"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/transportlayer/rest/handler/v2"
	"go-playground/cmd/api/internal/transportlayer/rest/oapi"
	"go-playground/pkg/errorx"
	"go-playground/pkg/ptr"
	"go-playground/pkg/timex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestListTasks(t *testing.T) {
	type input struct {
		w     *httptest.ResponseRecorder
		r     *http.Request
		param oapi.ListTasksParams
	}
	type want struct {
		status int
		body   string
	}
	tests := map[string]struct {
		input input
		setup func() *handler.TaskHandler
		want  want
	}{
		"success: no param": {
			input: input{
				w: httptest.NewRecorder(),
				r: httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/tasks", nil),
			},
			setup: func() *handler.TaskHandler {
				mck := new(MockTaskInteractor)
				mck.On("ListTasks", context.Background(), "", int32(0)).Return(entity.Page[entity.Task]{
					HasNext:   true,
					NextToken: "eyJpZCI6IjAxOTJiODQzLTE1MWUtNzRmZS04MTk4LTBlNjljZTM3OTMyYiJ9",
					Items: []entity.Task{
						{
							ID:        "0192b845-7a32-706b-ae58-d46437963c0e",
							Content:   "this is test",
							CreatedAt: time.Date(2024, 10, 23, 16, 26, 54, 0, timex.JST()),
							UpdatedAt: time.Date(2024, 10, 23, 16, 26, 54, 0, timex.JST()),
						},
					},
				}, nil)
				return &handler.TaskHandler{TaskInteractor: mck}
			},
			want: want{
				status: http.StatusOK,
				body: `
{
  "hasNext": true,
  "items": [
    {
      "content": "this is test",
      "createdAt": "2024-10-23T16:26:54+09:00",
      "id": "0192b845-7a32-706b-ae58-d46437963c0e",
      "updatedAt": "2024-10-23T16:26:54+09:00"
    }
  ],
  "next": "eyJpZCI6IjAxOTJiODQzLTE1MWUtNzRmZS04MTk4LTBlNjljZTM3OTMyYiJ9"
}
				`,
			},
		},
		"success: with next param": {
			input: input{
				w:     httptest.NewRecorder(),
				r:     httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/tasks?next=eyJpZCI6IjAxOTJiODQzLTE1MWUtNzRmZS04MTk4LTBlNjljZTM3OTMyYiJ9", nil),
				param: oapi.ListTasksParams{Next: ptr.String("eyJpZCI6IjAxOTJiODQzLTE1MWUtNzRmZS04MTk4LTBlNjljZTM3OTMyYiJ9")},
			},
			setup: func() *handler.TaskHandler {
				mck := new(MockTaskInteractor)
				mck.On("ListTasks", context.Background(), "eyJpZCI6IjAxOTJiODQzLTE1MWUtNzRmZS04MTk4LTBlNjljZTM3OTMyYiJ9", int32(0)).Return(entity.Page[entity.Task]{
					NextToken: "eyJpZCI6IjAxOTJiODNmLWUxOTktNzlkMS1hODcyLWIzZGNmMWY0MTE5YSJ9",
					HasNext:   true,
					Items: []entity.Task{
						{
							ID:        "0192b843-151e-74fe-8198-0e69ce37932b",
							Content:   "this is test 2",
							CreatedAt: time.Date(2024, 10, 23, 16, 24, 17, 0, timex.JST()),
							UpdatedAt: time.Date(2024, 10, 23, 16, 24, 17, 0, timex.JST()),
						},
					},
				}, nil)
				return &handler.TaskHandler{TaskInteractor: mck}
			},
			want: want{
				status: http.StatusOK,
				body: `
{
  "hasNext": true,
  "items": [
    {
      "content": "this is test 2",
      "createdAt": "2024-10-23T16:24:17+09:00",
      "id": "0192b843-151e-74fe-8198-0e69ce37932b",
      "updatedAt": "2024-10-23T16:24:17+09:00"
    }
  ],
  "next": "eyJpZCI6IjAxOTJiODNmLWUxOTktNzlkMS1hODcyLWIzZGNmMWY0MTE5YSJ9"
}
				`,
			},
		},
		"success: with limit param": {
			input: input{
				w:     httptest.NewRecorder(),
				r:     httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/tasks?limit=1", nil),
				param: oapi.ListTasksParams{Limit: ptr.Int32(1)},
			},
			setup: func() *handler.TaskHandler {
				mck := new(MockTaskInteractor)
				mck.On("ListTasks", context.Background(), "", int32(1)).Return(entity.Page[entity.Task]{
					NextToken: "eyJpZCI6IjAxOTJiODNmLWUxOTktNzlkMS1hODcyLWIzZGNmMWY0MTE5YSJ9",
					HasNext:   true,
					Items: []entity.Task{
						{
							ID:        "0192b843-151e-74fe-8198-0e69ce37932b",
							Content:   "this is test 2",
							CreatedAt: time.Date(2024, 10, 23, 16, 24, 17, 0, timex.JST()),
							UpdatedAt: time.Date(2024, 10, 23, 16, 24, 17, 0, timex.JST()),
						},
					},
				}, nil)
				return &handler.TaskHandler{TaskInteractor: mck}
			},
			want: want{
				status: http.StatusOK,
				body: `
{
  "hasNext": true,
  "items": [
    {
      "content": "this is test 2",
      "createdAt": "2024-10-23T16:24:17+09:00",
      "id": "0192b843-151e-74fe-8198-0e69ce37932b",
      "updatedAt": "2024-10-23T16:24:17+09:00"
    }
  ],
  "next": "eyJpZCI6IjAxOTJiODNmLWUxOTktNzlkMS1hODcyLWIzZGNmMWY0MTE5YSJ9"
}
				`,
			},
		},
		"success: no result": {
			input: input{
				w: httptest.NewRecorder(),
				r: httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/tasks", nil),
			},
			setup: func() *handler.TaskHandler {
				mck := new(MockTaskInteractor)
				mck.On("ListTasks", context.Background(), "", int32(0)).Return(entity.Page[entity.Task]{}, nil)
				return &handler.TaskHandler{TaskInteractor: mck}
			},
			want: want{
				status: http.StatusOK,
				body:   `{"hasNext":false,"next":"","items":[]}`,
			},
		},
		"failed to list tasks": {
			input: input{
				w: httptest.NewRecorder(),
				r: httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/tasks", nil),
			},
			setup: func() *handler.TaskHandler {
				mck := new(MockTaskInteractor)
				mck.On("ListTasks", context.Background(), "", int32(0)).Return(entity.Page[entity.Task]{}, errorx.NewWarn("failed to list task", errorx.WithStatus(http.StatusBadRequest)))
				return &handler.TaskHandler{TaskInteractor: mck}
			},
			want: want{
				status: http.StatusBadRequest,
				body:   `{"message":"failed to list task"}`,
			},
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			hn := v.setup()

			hn.ListTasks(v.input.w, v.input.r, v.input.param)

			assert.Equal(t, v.want.status, v.input.w.Code)
			assert.JSONEq(t, v.want.body, v.input.w.Body.String())
		})
	}
}

func TestGetTas(t *testing.T) {
	type input struct {
		w   *httptest.ResponseRecorder
		r   *http.Request
		tid oapi.TaskID
	}
	type want struct {
		status int
		body   string
	}
	tests := map[string]struct {
		input input
		setup func() *handler.TaskHandler
		want  want
	}{
		"success": {
			input: input{
				w:   httptest.NewRecorder(),
				r:   httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/tasks/0192b83f-e199-79d1-a872-b3dcf1f4119a", nil),
				tid: "0192b83f-e199-79d1-a872-b3dcf1f4119a",
			},
			setup: func() *handler.TaskHandler {
				mck := new(MockTaskInteractor)
				mck.On("FindTaskByID", context.Background(), "0192b83f-e199-79d1-a872-b3dcf1f4119a").Return(entity.Task{
					ID:        "0192b83f-e199-79d1-a872-b3dcf1f4119a",
					Content:   "this is test",
					CreatedAt: time.Date(2024, 10, 23, 16, 20, 47, 0, timex.JST()),
					UpdatedAt: time.Date(2024, 10, 23, 16, 20, 47, 0, timex.JST()),
				}, nil)
				return &handler.TaskHandler{TaskInteractor: mck}
			},
			want: want{
				status: http.StatusOK,
				body: `
{
  "content": "this is test",
  "createdAt": "2024-10-23T16:20:47+09:00",
  "id": "0192b83f-e199-79d1-a872-b3dcf1f4119a",
  "updatedAt": "2024-10-23T16:20:47+09:00"
}
				`,
			},
		},
		"failure": {
			input: input{
				w:   httptest.NewRecorder(),
				r:   httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/tasks/abc", nil),
				tid: "abc",
			},
			setup: func() *handler.TaskHandler {
				mck := new(MockTaskInteractor)
				mck.On("FindTaskByID", context.Background(), "abc").Return(entity.Task{}, errorx.NewWarn("missing task by abc", errorx.WithStatus(http.StatusNotFound)))
				return &handler.TaskHandler{TaskInteractor: mck}
			},
			want: want{
				status: http.StatusNotFound,
				body:   `{"message":"missing task by abc"}`,
			},
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			hn := v.setup()

			hn.GetTask(v.input.w, v.input.r, v.input.tid)

			assert.Equal(t, v.want.status, v.input.w.Code)
			assert.JSONEq(t, v.want.body, v.input.w.Body.String())
		})
	}
}

func TestPostTask(t *testing.T) {
	type input struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	type want struct {
		status int
		body   string
	}
	tests := map[string]struct {
		input input
		setup func() *handler.TaskHandler
		want  want
	}{
		"success": {
			input: input{
				w: httptest.NewRecorder(),
				r: httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/tasks", strings.NewReader(`{"content":"ok"}`)),
			},
			setup: func() *handler.TaskHandler {
				mck := new(MockTaskInteractor)
				mck.On("CreateTask", context.Background(), "ok").Return("0192b845-7a32-706b-ae58-d46437963c0e", nil)
				return &handler.TaskHandler{TaskInteractor: mck}
			},
			want: want{
				status: http.StatusOK,
				body:   `{"id":"0192b845-7a32-706b-ae58-d46437963c0e"}`,
			},
		},
		"failure: failed to unmarshal body": {
			input: input{
				w: httptest.NewRecorder(),
				r: httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/tasks", strings.NewReader(``)),
			},
			setup: func() *handler.TaskHandler { return &handler.TaskHandler{} },
			want: want{
				status: http.StatusBadRequest,
				body:   `{"message":"failed to unmarshal request body of PostTask"}`,
			},
		},
		"failure: failed to create task": {
			input: input{
				w: httptest.NewRecorder(),
				r: httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/tasks", strings.NewReader(`{"content":"failed"}`)),
			},
			setup: func() *handler.TaskHandler {
				mck := new(MockTaskInteractor)
				mck.On("CreateTask", context.Background(), "failed").Return("", errorx.NewError("failed to create task"))
				return &handler.TaskHandler{TaskInteractor: mck}
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   `{"message":"failed to create task"}`,
			},
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			hn := v.setup()

			hn.PostTask(v.input.w, v.input.r)

			assert.Equal(t, v.want.status, v.input.w.Code)
			assert.JSONEq(t, v.want.body, v.input.w.Body.String())
		})
	}
}

func TestPutTask(t *testing.T) {
	type input struct {
		w   *httptest.ResponseRecorder
		r   *http.Request
		tid oapi.TaskID
	}
	type want struct {
		status int
		body   string
	}
	tests := map[string]struct {
		input input
		setup func() *handler.TaskHandler
		want  want
	}{
		"success": {
			input: input{
				w:   httptest.NewRecorder(),
				r:   httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/tasks/0192b845-7a32-706b-ae58-d46437963c0e", strings.NewReader(`{"content":"want modify"}`)),
				tid: "0192b845-7a32-706b-ae58-d46437963c0e",
			},
			setup: func() *handler.TaskHandler {
				mck := new(MockTaskInteractor)
				mck.On("UpdateTask", context.Background(), "0192b845-7a32-706b-ae58-d46437963c0e", "want modify").Return(nil)
				return &handler.TaskHandler{TaskInteractor: mck}
			},
			want: want{
				status: http.StatusOK,
				body:   `{"id":"0192b845-7a32-706b-ae58-d46437963c0e"}`,
			},
		},
		"failure: failed to unmarshal body": {
			input: input{
				w:   httptest.NewRecorder(),
				r:   httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/tasks/0192b845-7a32-706b-ae58-d46437963c0e", strings.NewReader(``)),
				tid: "0192b845-7a32-706b-ae58-d46437963c0e",
			},
			setup: func() *handler.TaskHandler { return &handler.TaskHandler{} },
			want: want{
				status: http.StatusBadRequest,
				body:   `{"message":"failed to unmarshal request body of PutTask"}`,
			},
		},
		"failure: failed to update task": {
			input: input{
				w:   httptest.NewRecorder(),
				r:   httptest.NewRequestWithContext(context.Background(), http.MethodPut, "/tasks/0192b845-7a32-706b-ae58-d46437963c0e", strings.NewReader(`{"content":"failed"}`)),
				tid: "0192b845-7a32-706b-ae58-d46437963c0e",
			},
			setup: func() *handler.TaskHandler {
				mck := new(MockTaskInteractor)
				mck.On("UpdateTask", context.Background(), "0192b845-7a32-706b-ae58-d46437963c0e", "failed").Return(errorx.NewError("failed to update task"))
				return &handler.TaskHandler{TaskInteractor: mck}
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   `{"message":"failed to update task"}`,
			},
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			hn := v.setup()

			hn.PutTask(v.input.w, v.input.r, v.input.tid)

			assert.Equal(t, v.want.status, v.input.w.Code)
			assert.JSONEq(t, v.want.body, v.input.w.Body.String())
		})
	}
}
