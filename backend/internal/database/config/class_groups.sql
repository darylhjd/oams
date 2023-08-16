-- name: UpsertClassGroups :batchone
-- Insert a class group into the database. If the class group already exists, do nothing.
INSERT INTO class_groups (class_id, name, class_type, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT DO NOTHING
RETURNING *;