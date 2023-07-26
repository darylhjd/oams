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

-- name: UpsertUsers :batchone
INSERT INTO users (id, name, email, role, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
ON CONFLICT (id)
    DO UPDATE SET name       = $2,
                  email      = $3,
                  role       = $4,
                  updated_at = NOW()
RETURNING *;