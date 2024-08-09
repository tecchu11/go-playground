package handler

import (
	"encoding/json"
	"go-playground/pkg/errorx"
	"go-playground/pkg/httpx"
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func ListTasks(taskInteractor TaskInteractor) http.Handler {
	return ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		defer newrelic.FromContext(ctx).StartSegment("handler/ListTasks").End()

		limit, err := httpx.QueryInt32(r, "limit")
		if err != nil {
			return errorx.NewWarn("limit must be number", errorx.WithCause(err), errorx.WithStatus(400))
		}
		next := r.URL.Query().Get("next")
		tasks, err := taskInteractor.ListTasks(ctx, next, limit)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(tasks)
	})
}

func FindTaskByID(taskInteractor TaskInteractor) http.Handler {
	return ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		defer newrelic.FromContext(ctx).StartSegment("handler/FindTaskByID").End()

		id := r.PathValue("id")
		task, err := taskInteractor.FindTaskByID(ctx, id)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(task)
	})
}

func PostTask(taskInteractor TaskInteractor) http.Handler {
	return ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		defer newrelic.FromContext(ctx).StartSegment("handler/PostTask").End()

		var body ReqPostTask
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return err
		}
		id, err := taskInteractor.CreateTask(ctx, body.Content)
		if err != nil {
			return err
		}
		res := ResPostTask{id}
		w.WriteHeader(http.StatusCreated)
		return json.NewEncoder(w).Encode(res)
	})
}

func PutTask(taskInteractor TaskInteractor) http.Handler {
	return ErrorHandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()
		defer newrelic.FromContext(ctx).StartSegment("handler/PutTask").End()

		id := r.PathValue("id")
		var body ReqPutTask
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			return err
		}
		err := taskInteractor.UpdateTask(ctx, id, body.Content)
		if err != nil {
			return err
		}
		return json.NewEncoder(w).Encode(nil)
	})
}
