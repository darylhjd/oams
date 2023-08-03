-- name: ListSessionEnrollments :many
SELECT *
FROM session_enrollments
ORDER BY session_id, user_id;

-- name: GetSessionEnrollmentsBySessionID :many
SELECT *
FROM session_enrollments
WHERE session_id = $1;

-- name: GetSessionEnrollmentsByUserID :many
SELECT *
FROM session_enrollments
WHERE user_id = $1;

-- name: CreateSessionEnrollments :batchone
INSERT INTO session_enrollments (session_id, user_id, attended, created_at)
VALUES ($1, $2, $3, NOW())
RETURNING *;