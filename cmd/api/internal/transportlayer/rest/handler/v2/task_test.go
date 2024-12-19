package handler_test

import (
	"context"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/transportlayer/rest/handler/v2"
	"go-playground/cmd/api/internal/transportlayer/rest/oapi"
	"go-playground/pkg/apperr"
	"go-playground/pkg/ptr"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskHandler_ListTasks(t *testing.T) {
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
							CreatedAt: time.Date(2024, 10, 23, 16, 26, 54, 0, time.UTC),
							UpdatedAt: time.Date(2024, 10, 23, 16, 26, 54, 0, time.UTC),
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
      "createdAt": "2024-10-23T16:26:54Z",
      "id": "0192b845-7a32-706b-ae58-d46437963c0e",
      "updatedAt": "2024-10-23T16:26:54Z"
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
							CreatedAt: time.Date(2024, 10, 23, 16, 24, 17, 0, time.UTC),
							UpdatedAt: time.Date(2024, 10, 23, 16, 24, 17, 0, time.UTC),
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
      "createdAt": "2024-10-23T16:24:17Z",
      "id": "0192b843-151e-74fe-8198-0e69ce37932b",
      "updatedAt": "2024-10-23T16:24:17Z"
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
							CreatedAt: time.Date(2024, 10, 23, 16, 24, 17, 0, time.UTC),
							UpdatedAt: time.Date(2024, 10, 23, 16, 24, 17, 0, time.UTC),
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
      "createdAt": "2024-10-23T16:24:17Z",
      "id": "0192b843-151e-74fe-8198-0e69ce37932b",
      "updatedAt": "2024-10-23T16:24:17Z"
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
				mck.On("ListTasks", context.Background(), "", int32(0)).Return(entity.Page[entity.Task]{}, apperr.New("failed to list task", "failed to list task", apperr.CodeInvalidArgument))
				return &handler.TaskHandler{TaskInteractor: mck}
			},
			want: want{
				status: http.StatusBadRequest,
				body:   `{"message":"failed to list task"}`,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			hn := tc.setup()

			hn.ListTasks(tc.input.w, tc.input.r, tc.input.param)

			assert.Equal(t, tc.want.status, tc.input.w.Code)
			assert.JSONEq(t, tc.want.body, tc.input.w.Body.String())
		})
	}
}

func TestTaskHandler_GetTask(t *testing.T) {
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
					CreatedAt: time.Date(2024, 10, 23, 16, 20, 47, 0, time.UTC),
					UpdatedAt: time.Date(2024, 10, 23, 16, 20, 47, 0, time.UTC),
				}, nil)
				return &handler.TaskHandler{TaskInteractor: mck}
			},
			want: want{
				status: http.StatusOK,
				body: `
{
  "content": "this is test",
  "createdAt": "2024-10-23T16:20:47Z",
  "id": "0192b83f-e199-79d1-a872-b3dcf1f4119a",
  "updatedAt": "2024-10-23T16:20:47Z"
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
				mck.On("FindTaskByID", context.Background(), "abc").Return(entity.Task{}, apperr.New("missing task", "missing task by abc", apperr.CodeNotFound))
				return &handler.TaskHandler{TaskInteractor: mck}
			},
			want: want{
				status: http.StatusNotFound,
				body:   `{"message":"missing task by abc"}`,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			hn := tc.setup()

			hn.GetTask(tc.input.w, tc.input.r, tc.input.tid)

			assert.Equal(t, tc.want.status, tc.input.w.Code)
			assert.JSONEq(t, tc.want.body, tc.input.w.Body.String())
		})
	}
}

func TestTaskHandler_PostTask(t *testing.T) {
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
				body:   `{"message":"invalid request"}`,
			},
		},
		"failure: failed to create task": {
			input: input{
				w: httptest.NewRecorder(),
				r: httptest.NewRequestWithContext(context.Background(), http.MethodPost, "/tasks", strings.NewReader(`{"content":"failed"}`)),
			},
			setup: func() *handler.TaskHandler {
				mck := new(MockTaskInteractor)
				mck.On("CreateTask", context.Background(), "failed").Return("", apperr.New("internal server error", "failed to create new task", apperr.CodeInternal))
				return &handler.TaskHandler{TaskInteractor: mck}
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   `{"message":"failed to create new task"}`,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			hn := tc.setup()

			hn.PostTask(tc.input.w, tc.input.r)

			assert.Equal(t, tc.want.status, tc.input.w.Code)
			assert.JSONEq(t, tc.want.body, tc.input.w.Body.String())
		})
	}
}

func TestTaskHandler_PutTask(t *testing.T) {
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
				body:   `{"message":"invalid request"}`,
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
				mck.On("UpdateTask", context.Background(), "0192b845-7a32-706b-ae58-d46437963c0e", "failed").Return(apperr.New("internal server error", "failed to update task"))
				return &handler.TaskHandler{TaskInteractor: mck}
			},
			want: want{
				status: http.StatusInternalServerError,
				body:   `{"message":"failed to update task"}`,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			hn := tc.setup()

			hn.PutTask(tc.input.w, tc.input.r, tc.input.tid)

			assert.Equal(t, tc.want.status, tc.input.w.Code)
			assert.JSONEq(t, tc.want.body, tc.input.w.Body.String())
		})
	}
}
