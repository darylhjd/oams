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

-- name: UpdateSessionEnrollment :one
UPDATE session_enrollments
SET session_id = COALESCE(sqlc.narg('session_id'), session_id),
    user_id    = COALESCE(sqlc.narg('user_id'), user_id),
    attended   = COALESCE(sqlc.narg('attended'), attended),
    updated_at =
        CASE
            WHEN (NOT (sqlc.narg('session_id')::BIGINT IS NULL AND
                       sqlc.narg('user_id')::BIGINT IS NULL AND
                       sqlc.narg('attended')::BOOLEAN IS NULL))
                AND
                 (COALESCE(sqlc.narg('session_id'), session_id) <> session_id OR
                  COALESCE(sqlc.narg('user_id'), user_id) <> user_id OR
                  COALESCE(sqlc.narg('attended'), attended) <> attended)
                THEN NOW()
            ELSE updated_at END
WHERE id = $1
RETURNING id, session_id, user_id, attended, updated_at;

-- name: DeleteSessionEnrollment :one
DELETE
FROM session_enrollments
WHERE id = $1
RETURNING *;

-- name: CreateSessionEnrollments :batchone
INSERT INTO session_enrollments (session_id, user_id, attended, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING *;