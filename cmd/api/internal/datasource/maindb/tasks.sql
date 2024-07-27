-- name: FindTask :one
SELECT
    *
FROM
    tasks
WHERE
    id = ?;

-- name: CreateTask :execresult
INSERT INTO
    tasks (id, content)
VALUES
    (?, ?);

-- name: UpdateTask :execresult
UPDATE
    tasks
SET
    content = ?
WHERE
    id = ?;
