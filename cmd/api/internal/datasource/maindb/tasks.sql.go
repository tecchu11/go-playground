// Code generated by sqlc. DO NOT EDIT.
// source: tasks.sql

package maindb

import (
	"context"
	"database/sql"
)

const createTask = `-- name: CreateTask :execresult
INSERT INTO tasks (id, content)
		VALUES(?, ?)
`

type CreateTaskParams struct {
	ID      string
	Content string
}

// CreateTask inserts given task.
func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createTask, arg.ID, arg.Content)
}

const findTask = `-- name: FindTask :one
SELECT
	id, content, created_at, updated_at
FROM
	tasks
WHERE
	id = ?
`

// FindTask finds task by given id.
func (q *Queries) FindTask(ctx context.Context, id string) (Task, error) {
	row := q.db.QueryRowContext(ctx, findTask, id)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listTasks = `-- name: ListTasks :many
SELECT
    id,
    content,
    created_at,
    updated_at
FROM
    tasks
WHERE
    '' = ? OR id <= ?
ORDER BY
    id DESC
LIMIT ?
`

type ListTasksParams struct {
	ID    string
	Limit int32
}

// ListTasks finds tasks by cursor pagination.
func (q *Queries) ListTasks(ctx context.Context, arg ListTasksParams) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, listTasks, arg.ID, arg.ID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Content,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTask = `-- name: UpdateTask :execresult
UPDATE
	tasks
SET
	content = ?
WHERE
	id = ?
`

type UpdateTaskParams struct {
	Content string
	ID      string
}

// UpdateTask updates task by given id.
func (q *Queries) UpdateTask(ctx context.Context, arg UpdateTaskParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, updateTask, arg.Content, arg.ID)
}
