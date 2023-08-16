-- name: UpdateSessionEnrollment :one
UPDATE session_enrollments
SET attended   = COALESCE(sqlc.narg('attended'), attended)
WHERE id = $1
RETURNING id, session_id, user_id, attended, updated_at;

-- name: UpsertSessionEnrollments :batchone
-- Insert a session enrollment into the database. If the session enrollment already exists, do nothing.
INSERT INTO session_enrollments (session_id, user_id, attended, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT DO NOTHING
RETURNING *;