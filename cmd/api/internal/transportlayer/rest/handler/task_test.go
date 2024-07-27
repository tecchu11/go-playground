package handler_test

import (
	"context"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/transportlayer/rest/handler"
	"go-playground/pkg/errorx"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFindTaskByID(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}
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
					CreatedAt: time.Date(2024, 7, 26, 0, 0, 0, 0, loc),
					UpdatedAt: time.Date(2024, 7, 26, 0, 0, 0, 0, loc),
				},
				nil,
			},
			expectedCode: 200,
			expectedBody: `{"ID":"test-id", "Content":"do test", "CreatedAt":"2024-07-26T00:00:00+09:00", "UpdatedAt":"2024-07-26T00:00:00+09:00"}`,
		},
		"not found task": {
			w:                httptest.NewRecorder(),
			r:                httptest.NewRequest("", "/tasks/{id}", nil),
			id:               "missing-id",
			mockFindTaskByID: []any{entity.Task{}, errorx.NewInfo("missing task", errorx.WithStatus(404))},
			expectedCode:     404,
			expectedBody:     `{"title":"Handled error", "type":"about:blank", "detail":"missing task", "instance":"/tasks/{id}", "status":404}`,
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
		mockCreateTask error
		expectedCode   int
		expectedBody   string
	}{
		"success": {
			w:              httptest.NewRecorder(),
			r:              httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"content":"do test"}`)),
			mockCreateTask: nil,
			expectedCode:   201,
			expectedBody:   `null`,
		},
		"failed to create task": {
			w:              httptest.NewRecorder(),
			r:              httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"content":"do test"}`)),
			mockCreateTask: errorx.NewError("failed to create task"),
			expectedCode:   500,
			expectedBody:   `{"type":"about:blank", "title":"Handled error", "detail":"failed to create task", "instance":"/tasks", "status":500}`,
		},
		"unmarshal error": {
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest("POST", "/tasks", strings.NewReader("")),
			expectedCode: 500,
			expectedBody: `{"type":"about:blank", "title":"Unhandled error", "detail":"EOF", "instance":"/tasks", "status":500}`,
		},
	}

	for k, v := range tests {
		t.Run(k, func(t *testing.T) {
			mockTakInteractor := MockTaskInteractor{}
			mockTakInteractor.On("CreateTask", context.Background(), "do test").Return(v.mockCreateTask)

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
			expectedBody:   `{"type":"about:blank", "title":"Handled error", "detail":"failed to update task", "instance":"/tasks/{id}", "status":500}`,
		},
		"unmarshal error": {
			w:            httptest.NewRecorder(),
			r:            httptest.NewRequest("PUT", "/tasks/{id}", strings.NewReader("")),
			expectedCode: 500,
			expectedBody: `{"type":"about:blank", "title":"Unhandled error", "detail":"EOF", "instance":"/tasks/{id}", "status":500}`,
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
