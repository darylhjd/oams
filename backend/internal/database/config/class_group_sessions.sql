-- name: UpdateClassGroupSession :one
UPDATE class_group_sessions
SET class_group_id = COALESCE(sqlc.narg('class_group_id'), class_group_id),
    start_time     = COALESCE(sqlc.narg('start_time'), start_time),
    end_time       = COALESCE(sqlc.narg('end_time'), end_time),
    venue          = COALESCE(sqlc.narg('venue'), venue),
    updated_at     =
        CASE
            WHEN (COALESCE(sqlc.narg('class_group_id'), class_group_id) <> class_group_id OR
                  COALESCE(sqlc.narg('start_time'), start_time) <> start_time OR
                  COALESCE(sqlc.narg('end_time'), end_time) <> end_time OR
                  COALESCE(sqlc.narg('venue'), venue) <> venue)
                THEN NOW()
            ELSE updated_at
            END
WHERE id = $1
RETURNING id, class_group_id, start_time, end_time, venue, updated_at;

-- name: UpsertClassGroupSessions :batchone
-- Insert a class group session into the database. If the class group session already exists, only update
-- the end_time and venue.
INSERT INTO class_group_sessions (class_group_id, start_time, end_time, venue, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
ON CONFLICT ON CONSTRAINT ux_class_group_id_start_time
    DO UPDATE SET end_time   = $3,
                  venue      = $4,
                  updated_at =
                      CASE
                          WHEN $3 <> class_group_sessions.end_time OR
                               $4 <> class_group_sessions.venue
                              THEN NOW()
                          ELSE class_group_sessions.updated_at
                          END
RETURNING *;