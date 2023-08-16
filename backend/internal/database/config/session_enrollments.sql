-- name: UpsertSessionEnrollments :batchone
-- Insert a session enrollment into the database. If the session enrollment already exists, do nothing.
INSERT INTO session_enrollments (session_id, user_id, attended, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT DO NOTHING
RETURNING *;