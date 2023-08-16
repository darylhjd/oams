-- name: UpsertUsers :batchone
-- Insert a user into the database. If the user already exists, then only update the name and email.
INSERT INTO users (id, name, email, role, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
ON CONFLICT (id)
    DO UPDATE SET name       = $2,
                  email      = $3
RETURNING *;