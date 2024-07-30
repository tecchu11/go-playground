-- name: ListTasks :many
-- ListTasks finds tasks by cursor pagination.
SELECT
	*
FROM
	tasks
WHERE
	id >= ?
ORDER BY
	id
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
