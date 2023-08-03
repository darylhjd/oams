-- name: ListSessionEnrollments :many
SELECT *
FROM session_enrollments
ORDER BY session_id, user_id;

-- name: GetSessionEnrollment :one
SELECT *
FROM session_enrollments
WHERE id = $1
LIMIT 1;

-- name: GetSessionEnrollmentsBySessionID :many
SELECT *
FROM session_enrollments
WHERE session_id = $1;

-- name: GetSessionEnrollmentsByUserID :many
SELECT *
FROM session_enrollments
WHERE user_id = $1;

-- name: CreateSessionEnrollment :one
INSERT INTO session_enrollments (session_id, user_id, attended, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING id, session_id, user_id, attended, created_at;

-- name: CreateSessionEnrollments :batchone
INSERT INTO session_enrollments (session_id, user_id, attended, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING *;