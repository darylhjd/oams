-- name: UpsertClassGroupSessions :batchone
-- Insert a class group session into the database. If the class group session already exists, only update
-- the end_time and venue.
INSERT INTO class_group_sessions (class_group_id, start_time, end_time, venue, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
ON CONFLICT ON CONSTRAINT ux_class_group_id_start_time
    DO UPDATE SET end_time   = $3,
                  venue      = $4
RETURNING *;