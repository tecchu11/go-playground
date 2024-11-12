-- name: CreateUser :execrows
-- CreateUser inserts given user.
INSERT INTO
	users (
		id,
		sub,
		given_name,
		family_name,
		email,
		email_verified
	)
VALUES
	(?, ?, ?, ?, ?, ?);

-- name: FindUserBySub :one
--  FindUserBySub finds user with given sub(jwt subject).
SELECT
	id,
	sub,
	given_name,
	family_name,
	email,
	email_verified,
	created_at,
	updated_at
FROM
	users
WHERE
	sub = ?;
