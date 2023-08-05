-- name: ListClassGroups :many
SELECT *
FROM class_groups
ORDER BY class_id, name;

-- name: GetClassGroup :one
SELECT *
FROM class_groups
WHERE id = $1
LIMIT 1;

-- name: GetClassGroupsByIDs :many
SELECT *
FROM class_groups
WHERE id = ANY (@ids::BIGINT[])
ORDER BY id;

-- name: CreateClassGroup :one
INSERT INTO class_groups (class_id, name, class_type, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING id, class_id, name, class_type, created_at;

-- name: UpdateClassGroup :one
UPDATE class_groups
SET class_id   = COALESCE(sqlc.narg('class_id'), class_id),
    name       = COALESCE(sqlc.narg('name'), name),
    class_type = COALESCE(sqlc.narg('class_type'), class_type),
    updated_at =
        CASE
            WHEN (COALESCE(sqlc.narg('class_id'), class_id) <> class_id OR
                  COALESCE(sqlc.narg('name'), name) <> name OR
                  COALESCE(sqlc.narg('class_type'), class_type) <> class_type)
                THEN NOW()
            ELSE updated_at END
WHERE id = $1
RETURNING id, class_id, name, class_type, updated_at;

-- name: DeleteClassGroup :one
DELETE
FROM class_groups
WHERE id = $1
RETURNING *;

-- name: UpsertClassGroups :batchone
-- Use this sparingly. When in doubt, use the atomic INSERT and UPDATE statements instead.
INSERT INTO class_groups (class_id, name, class_type, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT ON CONSTRAINT ux_class_id_name
    DO UPDATE SET class_type = $3,
                  updated_at =
                      CASE
                          WHEN $3 <> class_groups.class_type
                              THEN NOW()
                          ELSE class_groups.updated_at
                          END
RETURNING *;