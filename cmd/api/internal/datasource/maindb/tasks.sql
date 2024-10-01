-- name: ListTasks :many
-- ListTasks finds tasks by cursor pagination.
SELECT
    id,
    content,
    created_at,
    updated_at
FROM
    tasks
WHERE
    '' = sqlc.arg('id') OR id <= sqlc.arg('id')
ORDER BY
    id DESC
LIMIT ?;


-- name: FindTask :one
-- FindTask finds task by given id.
SELECT
	*
FROM
	tasks
WHERE
	id = ?;

-- name: CreateTask :execresult
-- CreateTask inserts given task.
INSERT INTO tasks (id, content)
		VALUES(?, ?);
 
-- name: UpdateTask :execresult
-- UpdateTask updates task by given id.
UPDATE
	tasks
SET
	content = ?
WHERE
	id = ?;
