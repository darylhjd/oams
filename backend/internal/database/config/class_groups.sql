-- name: UpdateClassGroup :one
UPDATE class_groups
SET class_id   = COALESCE(sqlc.narg('class_id'), class_id),
    name       = COALESCE(sqlc.narg('name'), name),
    class_type = COALESCE(sqlc.narg('class_type'), class_type)
WHERE id = $1
RETURNING id, class_id, name, class_type, updated_at;

-- name: UpsertClassGroups :batchone
-- Insert a class group into the database. If the class group already exists, do nothing.
INSERT INTO class_groups (class_id, name, class_type, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT DO NOTHING
RETURNING *;