-- name: ListSessionEnrollments :many
SELECT *
FROM session_enrollments
ORDER BY session_id, user_id;

-- name: GetSessionEnrollment :one
SELECT *
FROM session_enrollments
WHERE id = $1
LIMIT 1;

-- name: CreateSessionEnrollment :one
INSERT INTO session_enrollments (session_id, user_id, attended, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING id, session_id, user_id, attended, created_at;

-- name: UpdateSessionEnrollment :one
UPDATE session_enrollments
SET attended   = COALESCE(sqlc.narg('attended'), attended),
    updated_at =
        CASE
            WHEN COALESCE(sqlc.narg('attended'), attended) <> attended
                THEN NOW()
            ELSE updated_at END
WHERE id = $1
RETURNING id, session_id, user_id, attended, updated_at;

-- name: DeleteSessionEnrollment :one
DELETE
FROM session_enrollments
WHERE id = $1
RETURNING *;

-- name: UpsertSessionEnrollments :batchone
-- Insert a session enrollment into the database. If the session enrollment already exists, do nothing.
INSERT INTO session_enrollments (session_id, user_id, attended, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT DO NOTHING
RETURNING *;