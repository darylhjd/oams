-- name: ListUsers :many
SELECT *
FROM users
ORDER BY id;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUsersByIDs :many
SELECT *
FROM users
WHERE id = ANY (sqlc.arg(ids)::TEXT[])
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (id, name, email, role, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id, name, email, role, created_at;

-- name: UpdateUser :one
UPDATE users
SET name       = COALESCE(sqlc.narg('name'), name),
    email      = COALESCE(sqlc.narg('email'), email),
    role       = COALESCE(sqlc.narg('role'), role),
    updated_at = CASE
                     WHEN (NOT (sqlc.narg('name')::TEXT IS NULL AND
                                sqlc.narg('email')::TEXT IS NULL AND
                                sqlc.narg('role')::USER_ROLE IS NULL)) AND
                          (COALESCE(sqlc.narg('name'), name) <> name OR
                           COALESCE(sqlc.narg('email'), email) <> email OR
                           COALESCE(sqlc.narg('role'), role) <> role) THEN NOW()
                     ELSE updated_at END
WHERE id = $1
RETURNING id, name, email, role, updated_at;

-- name: DeleteUser :one
DELETE
FROM users
WHERE id = $1
RETURNING *;

-- name: UpsertUsers :batchone
INSERT INTO users (id, name, email, role, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
ON CONFLICT (id)
    DO UPDATE SET name       = $2,
                  email      = $3,
                  role       = $4,
                  updated_at = NOW()
RETURNING *;