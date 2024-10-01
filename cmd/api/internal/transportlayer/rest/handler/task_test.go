package handler_test

import (
	"context"
	"errors"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/transportlayer/rest/handler"
	"go-playground/pkg/errorx"
	"go-playground/pkg/timex"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestListTasks(t *testing.T) {
	tests := map[string]struct {
		w            *httptest.ResponseRecorder
		r            *http.Request
		mockOn       []any
		mockReturn   []any
		expectedCode int
		expectedBody string
	}{
		"success": {
			w:      httptest.NewRecorder(),
			r:      httptest.NewRequest("", "/tasks?limit=2", nil),
			mockOn: []any{context.Background(), "", int32(2)},
			mockReturn: []any{entity.Page[entity.Task]{
				Items:     []entity.Task{{ID: "test-id-1"}, {ID: "test-id-2"}},
				HasNext:   true,
				NextToken: "test-id-3",
			}, nil},
			expectedCode: 200,
			expectedBody: `{"items":[{"id":"test-id-1", "content":"", "createdAt":"0001-01-01T00:00:00Z", "updatedAt":"0001-01-01T00:00:00Z"},{"id":"test-id-2", "content":"", "createdAt":"0001-01-01T00:00:00Z", "updatedAt":"0001-01-01T00:00:00Z"}], "next":"test-id-3", "hasNext":true}`,
		},
		"limit is not number": {
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest("", "/tasks?limit=o", nil),
			expectedCode: 400,
			expectedBody: `{"message":"limit must be number"}`,
		},
		"failed to list tasks": {
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest("", "/tasks?limit=100&next=test-id", nil),
			mockOn:       []any{context.Background(), "test-id", int32(100)},
			mockReturn:   []any{entity.Page[entity.Task]{}, errors.New("unknown error on list tasks")},
			expectedCode: 500,
			expectedBody: `{"message":"unknown error on list tasks"}`,
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			mockTaskInteractor := MockTaskInteractor{}
			mockTaskInteractor.On("ListTasks", v.mockOn...).Return(v.mockReturn...)

			handler.ListTasks(&mockTaskInteractor).ServeHTTP(v.w, v.r)

			assert.Equal(t, v.expectedCode, v.w.Code)
			assert.JSONEq(t, v.expectedBody, v.w.Body.String())
		})
	}
}

func TestFindTaskByID(t *testing.T) {
	tests := map[string]struct {
		w                *httptest.ResponseRecorder
		r                *http.Request
		id               string
		mockFindTaskByID []any
		expectedCode     int
		expectedBody     string
	}{
		"success": {
			w:  httptest.NewRecorder(),
			r:  httptest.NewRequest("", "/tasks/{id}", nil),
			id: "test-id",
			mockFindTaskByID: []any{
				entity.Task{
					ID:        "test-id",
					Content:   "do test",
					CreatedAt: time.Date(2024, 7, 26, 0, 0, 0, 0, timex.JST()),
					UpdatedAt: time.Date(2024, 7, 26, 0, 0, 0, 0, timex.JST()),
				},
				nil,
			},
			expectedCode: 200,
			expectedBody: `{"id":"test-id", "content":"do test", "createdAt":"2024-07-26T00:00:00+09:00", "updatedAt":"2024-07-26T00:00:00+09:00"}`,
		},
		"not found task": {
			w:                httptest.NewRecorder(),
			r:                httptest.NewRequest("", "/tasks/{id}", nil),
			id:               "missing-id",
			mockFindTaskByID: []any{entity.Task{}, errorx.NewInfo("missing task", errorx.WithStatus(404))},
			expectedCode:     404,
			expectedBody:     `{"message":"missing task"}`,
		},
	}
	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			mockTaskInteractor := MockTaskInteractor{}
			mockTaskInteractor.On("FindTaskByID", context.Background(), v.id).Return(v.mockFindTaskByID...)
			v.r.SetPathValue("id", v.id)

			handler.FindTaskByID(&mockTaskInteractor).ServeHTTP(v.w, v.r)

			assert.Equal(t, v.expectedCode, v.w.Code)
			assert.JSONEq(t, v.expectedBody, v.w.Body.String())
		})
	}
}

func TestPostTask(t *testing.T) {
	tests := map[string]struct {
		w              *httptest.ResponseRecorder
		r              *http.Request
		mockCreateTask []any
		expectedCode   int
		expectedBody   string
	}{
		"success": {
			w:              httptest.NewRecorder(),
			r:              httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"content":"do test"}`)),
			mockCreateTask: []any{"019102c9-19e3-76cd-85cc-08a12fcfa8f9", nil},
			expectedCode:   201,
			expectedBody:   `{"id":"019102c9-19e3-76cd-85cc-08a12fcfa8f9"}`,
		},
		"failed to create task": {
			w:              httptest.NewRecorder(),
			r:              httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"content":"do test"}`)),
			mockCreateTask: []any{"", errorx.NewError("failed to create task")},
			expectedCode:   500,
			expectedBody:   `{"message":"failed to create task"}`,
		},
		"unmarshal error": {
			w:              httptest.NewRecorder(),
			r:              httptest.NewRequest("POST", "/tasks", strings.NewReader("")),
			mockCreateTask: []any{nil, nil},
			expectedCode:   500,
			expectedBody:   `{"message":"EOF"}`,
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			mockTakInteractor := MockTaskInteractor{}
			mockTakInteractor.On("CreateTask", context.Background(), "do test").Return(v.mockCreateTask[0], v.mockCreateTask[1])

			handler.PostTask(&mockTakInteractor).ServeHTTP(v.w, v.r)

			assert.Equal(t, v.expectedCode, v.w.Code)
			assert.JSONEq(t, v.expectedBody, v.w.Body.String())
		})
	}
}

func TestPutTask(t *testing.T) {
	tests := map[string]struct {
		w              *httptest.ResponseRecorder
		r              *http.Request
		id             string
		mockUpdateTask error
		expectedCode   int
		expectedBody   string
	}{
		"success": {
			w:              httptest.NewRecorder(),
			r:              httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"content":"do test"}`)),
			mockUpdateTask: nil,
			expectedCode:   200,
			expectedBody:   `null`,
		},
		"failed to create task": {
			w:              httptest.NewRecorder(),
			r:              httptest.NewRequest("PUT", "/tasks/{id}", strings.NewReader(`{"content":"do test"}`)),
			mockUpdateTask: errorx.NewError("failed to update task"),
			expectedCode:   500,
			expectedBody:   `{"message":"failed to update task"}`,
		},
		"unmarshal error": {
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest("PUT", "/tasks/{id}", strings.NewReader("")),
			expectedCode: 500,
			expectedBody: `{"message":"EOF"}`,
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			v.r.SetPathValue("id", v.id)
			mockTakInteractor := MockTaskInteractor{}
			mockTakInteractor.On("UpdateTask", context.Background(), v.id, "do test").Return(v.mockUpdateTask)

			handler.PutTask(&mockTakInteractor).ServeHTTP(v.w, v.r)

			assert.Equal(t, v.expectedCode, v.w.Code)
			assert.JSONEq(t, v.expectedBody, v.w.Body.String())
		})
	}
}
