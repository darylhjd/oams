-- name: CreateUser :one
INSERT INTO users (id, name, email, role, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id, name, email, role, created_at;

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

-- name: DeleteUser :one
DELETE
FROM users
WHERE id = $1
RETURNING *;

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

-- name: GetUserUpcomingClassGroupSessions :many
-- Get information on a user's upcoming classes. This query returns all session enrollments for that user that are
-- currently happening or will happen in the future. The sessions are returned in ascending order of start time and then
-- end time.
SELECT c.code,
       c.year,
       c.semester,
       cg.name,
       cg.class_type,
       cgs.start_time,
       cgs.end_time,
       cgs.venue
FROM class_group_sessions cgs
         INNER JOIN class_groups cg
                    ON cgs.class_group_id = cg.id
         INNER JOIN classes c
                    ON cg.class_id = c.id
WHERE cgs.id IN (SELECT session_id
                 FROM session_enrollments
                 WHERE user_id = $1)
  AND cgs.end_time > NOW()
ORDER BY cgs.start_time, cgs.end_time;