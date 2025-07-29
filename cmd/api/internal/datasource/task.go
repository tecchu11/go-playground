package datasource

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-playground/cmd/api/internal/datasource/database"
	"go-playground/cmd/api/internal/domain/entity"
	"go-playground/cmd/api/internal/domain/repository"
	"go-playground/pkg/apperr"

	"github.com/jmoiron/sqlx"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// TaskAdaptor is implementation of repository.TaskRepository.
type TaskAdaptor struct {
	base
}

// NewTaskAdaptor initializes TaskAdaptor.
func NewTaskAdaptor(db *sqlx.DB) *TaskAdaptor {
	return &TaskAdaptor{base: base{db: db}}
}

// ListTasks list all task.
func (a *TaskAdaptor) ListTasks(ctx context.Context, next entity.TaskID, limit int32) (entity.Page[entity.Task], error) {
	defer newrelic.FromContext(ctx).StartSegment("datasource/TaskAdaptor/ListTasks").End()

	queries := a.queriesFromContext(ctx)
	rows, err := queries.ListTasks(ctx, database.ListTasksParams{ID: next, Limit: limit + 1})
	if err != nil {
		return entity.Page[entity.Task]{}, apperr.New("list tasks", "failed to list tasks", apperr.WithCause(err))
	}
	tasks := make([]entity.Task, len(rows))
	for i, r := range rows {
		tasks[i] = entity.Task{
			ID:        r.ID,
			Content:   r.Content,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		}
	}
	return entity.NewPage(tasks, limit)
}

// FindByID select task from task record by given id. Error will be returned task is not found.
func (a *TaskAdaptor) FindByID(ctx context.Context, id entity.TaskID) (entity.Task, error) {
	defer newrelic.FromContext(ctx).StartSegment("datasource/TaskAdaptor/FindByID").End()

	queries := a.queriesFromContext(ctx)
	row, err := queries.FindTask(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Task{}, apperr.New(fmt.Sprintf("find task by id %q", id), "not found task", apperr.WithCause(err), apperr.CodeNotFound)
		}
		return entity.Task{}, apperr.New("find task", "failed to find task", apperr.WithCause(err))
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

	queries := a.queriesFromContext(ctx)
	_, err := queries.CreateTask(ctx, database.CreateTaskParams{
		ID:      task.ID,
		Content: task.Content,
	})
	if err != nil {
		return apperr.New("create new task", "failed to create new task", apperr.WithCause(err))
	}
	return nil
}

// Update task record by give task entity.
func (a *TaskAdaptor) Update(ctx context.Context, task entity.Task) error {
	defer newrelic.FromContext(ctx).StartSegment("datasource/TaskAdaptor/Update").End()

	queries := a.queriesFromContext(ctx)
	_, err := queries.UpdateTask(ctx, database.UpdateTaskParams{
		ID:      task.ID,
		Content: task.Content,
	})
	if err != nil {
		return apperr.New(fmt.Sprintf("update task by id %q", task.ID), "failed to update task", apperr.WithCause(err))
	}
	return nil
}

// Creates creates multiple tasks.
func (a *TaskAdaptor) Creates(ctx context.Context, tasks []entity.Task) error {
	ext := a.extFromContext(ctx)

	_, err := sqlx.NamedExec(ext, `INSERT INTO tasks (id, content) VALUES(:id, :content)`, tasks)
	if err != nil {
		return fmt.Errorf("named exec on creates: %w", err)
	}
	return nil
}

var _ repository.TaskRepository = (*TaskAdaptor)(nil)
