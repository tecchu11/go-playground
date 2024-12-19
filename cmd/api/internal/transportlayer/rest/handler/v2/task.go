package handler

import (
	"encoding/json"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/transportlayer/rest/oapi"
	"go-playground/pkg/apperr"
	"go-playground/pkg/collection"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type TaskHandler struct {
	TaskInteractor TaskInteractor
}

// ListTask lists tasks for [GET /tasks]
func (t *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request, params oapi.ListTasksParams) {
	defer newrelic.FromContext(r.Context()).StartSegment("handler/taskHandler/ListTasks").End()

	ErrorHandlerFunc(w, r, func(w http.ResponseWriter, r *http.Request) error {
		var (
			next  string
			limit int32
		)
		if params.Next != nil {
			next = *params.Next
		}
		if params.Limit != nil {
			limit = *params.Limit
		}
		result, err := t.TaskInteractor.ListTasks(r.Context(), next, limit)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(
			oapi.ResponseTasks{
				Next:    result.NextToken,
				HasNext: result.HasNext,
				Items: collection.SMap(result.Items, func(e entity.Task) oapi.Task {
					return oapi.Task{
						ID:        e.ID,
						Content:   e.Content,
						CreatedAt: e.CreatedAt,
						UpdatedAt: e.UpdatedAt,
					}
				}),
			},
		)
	})
}

// GetTask gets task by given id for [GET /tasks/{taskId}]
func (t *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request, id oapi.TaskID) {
	defer newrelic.FromContext(r.Context()).StartSegment("/handler/taskHandler/GetTask").End()

	ErrorHandlerFunc(w, r, func(w http.ResponseWriter, r *http.Request) error {
		result, err := t.TaskInteractor.FindTaskByID(r.Context(), id)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(
			oapi.ResponseTask{
				ID:        result.ID,
				Content:   result.Content,
				CreatedAt: result.CreatedAt,
				UpdatedAt: result.UpdatedAt,
			},
		)
	})
}

// PostTask post task with give content for [POST /tasks]
func (t *TaskHandler) PostTask(w http.ResponseWriter, r *http.Request) {
	defer newrelic.FromContext(r.Context()).StartSegment("handler/taskHandler/PostTask")

	ErrorHandlerFunc(w, r, func(w http.ResponseWriter, r *http.Request) error {
		var body oapi.PostTaskJSONRequestBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			return apperr.New("unmarshal PostTask", "invalid request", apperr.WithCause(err), apperr.CodeInvalidArgument)
		}
		result, err := t.TaskInteractor.CreateTask(r.Context(), body.Content)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(oapi.ResponseTaskID{
			ID: result,
		})
	})
}

// PutTask put task with given content by id for [PUT /tasks/{taskId}]
func (t *TaskHandler) PutTask(w http.ResponseWriter, r *http.Request, id oapi.TaskID) {
	defer newrelic.FromContext(r.Context()).StartSegment("handler/taskHandler/PutTask")

	ErrorHandlerFunc(w, r, func(w http.ResponseWriter, r *http.Request) error {
		var body oapi.PutTaskJSONRequestBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			return apperr.New("unmarshal PutTask body", "invalid request", apperr.WithCause(err), apperr.CodeInvalidArgument)
		}
		err = t.TaskInteractor.UpdateTask(r.Context(), id, body.Content)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(oapi.ResponseTaskID{
			ID: id,
		})
	})
}
