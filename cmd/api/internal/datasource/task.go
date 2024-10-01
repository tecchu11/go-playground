package datasource

import (
	"context"
	"database/sql"
	"errors"
	"go-playground/cmd/api/internal/datasource/maindb"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/domain/repository"
	"go-playground/pkg/errorx"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// TaskAdaptor is implementation of repository.TaskRepository.
type TaskAdaptor struct {
	queries *maindb.Queries
}

// NewTaskAdaptor initializes TaskAdaptor.
func NewTaskAdaptor(queries *maindb.Queries) *TaskAdaptor {
	return &TaskAdaptor{queries}
}

// ListTasks list all task.
func (a *TaskAdaptor) ListTasks(ctx context.Context, next entity.TaskID, limit int32) (entity.Page[entity.Task], error) {
	defer newrelic.FromContext(ctx).StartSegment("datasource/TaskAdaptor/ListTasks").End()

	queries := TransactionQueries(ctx, a.queries)
	rows, err := queries.ListTasks(ctx, maindb.ListTasksParams{ID: next, Limit: limit + 1})
	if err != nil {
		return entity.Page[entity.Task]{}, errorx.NewError("cant list tasks", errorx.WithCause(err))
	}
	var tasks []entity.Task
	for _, record := range rows {
		tasks = append(tasks, entity.Task{
			ID:        record.ID,
			Content:   record.Content,
			CreatedAt: record.CreatedAt,
			UpdatedAt: record.UpdatedAt,
		})
	}
	return entity.NewPage(tasks, limit)
}

// FindByID select task from task record by given id. Error will be returned task is not found.
func (a *TaskAdaptor) FindByID(ctx context.Context, id entity.TaskID) (entity.Task, error) {
	defer newrelic.FromContext(ctx).StartSegment("datasource/TaskAdaptor/FindByID").End()

	queries := TransactionQueries(ctx, a.queries)
	row, err := queries.FindTask(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Task{}, errorx.NewInfo(
				"not found task by id",
				errorx.WithStatus(404),
				errorx.WithCause(err),
			)
		}
		return entity.Task{}, errorx.NewError("Failed to find task by id", errorx.WithCause(err))
	}
	return entity.Task{
		ID:        row.ID,
		Content:   row.Content,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

// Create inserts given task to task table.
func (a *TaskAdaptor) Create(ctx context.Context, task entity.Task) error {
	defer newrelic.FromContext(ctx).StartSegment("datasource/TaskAdaptor/Create").End()

	queries := TransactionQueries(ctx, a.queries)
	result, err := queries.CreateTask(ctx, maindb.CreateTaskParams{
		ID:      task.ID,
		Content: task.Content,
	})
	if err != nil {
		return errorx.NewError("Failed to create task", errorx.WithCause(err))
	}
	_, err = result.RowsAffected()
	if err != nil {
		return errorx.NewError("Failed to create task", errorx.WithCause(err))
	}
	return nil
}

// Update task record by give task entity.
func (a *TaskAdaptor) Update(ctx context.Context, task entity.Task) error {
	defer newrelic.FromContext(ctx).StartSegment("datasource/TaskAdaptor/Update").End()

	queries := TransactionQueries(ctx, a.queries)
	result, err := queries.UpdateTask(ctx, maindb.UpdateTaskParams{
		ID:      task.ID,
		Content: task.Content,
	})
	if err != nil {
		return errorx.NewError("Failed to update task")
	}
	_, err = result.RowsAffected()
	if err != nil {
		return errorx.NewError("Failed to update task", errorx.WithCause(err))
	}
	return nil
}

var _ repository.TaskRepository = (*TaskAdaptor)(nil)
