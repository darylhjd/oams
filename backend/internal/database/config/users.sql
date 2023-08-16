-- name: UpdateUser :one
UPDATE users
SET name       = COALESCE(sqlc.narg('name'), name),
    email      = COALESCE(sqlc.narg('email'), email),
    role       = COALESCE(sqlc.narg('role'), role),
    updated_at =
        CASE
            WHEN (COALESCE(sqlc.narg('name'), name) <> name OR
                  COALESCE(sqlc.narg('email'), email) <> email OR
                  COALESCE(sqlc.narg('role'), role) <> role)
                THEN NOW()
            ELSE updated_at
            END
WHERE id = $1
RETURNING id, name, email, role, updated_at;

-- name: UpsertUsers :batchone
-- Insert a user into the database. If the user already exists, then only update the name and email.
INSERT INTO users (id, name, email, role, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
ON CONFLICT (id)
    DO UPDATE SET name       = $2,
                  email      = $3,
                  updated_at =
                      CASE
                          WHEN $2 <> users.name OR
                               $3 <> users.email OR
                               $4 <> users.role
                              THEN NOW()
                          ELSE users.updated_at
                          END
RETURNING *;