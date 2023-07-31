-- name: ListClasses :many
SELECT *
FROM classes
ORDER BY code, year, semester;

-- name: GetClass :one
SELECT *
FROM classes
WHERE id = $1
LIMIT 1;

-- name: GetClassesByIDs :many
SELECT *
FROM classes
WHERE id = ANY (@ids::BIGINT[])
ORDER BY id;

-- name: CreateClass :one
INSERT INTO classes (code, year, semester, programme, au, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
RETURNING id, code, year, semester, programme, au, created_at;

-- name: UpdateClass :one
UPDATE classes
SET code       = COALESCE(sqlc.narg('code'), code),
    year       = COALESCE(sqlc.narg('year'), year),
    semester   = COALESCE(sqlc.narg('semester'), semester),
    programme  = COALESCE(sqlc.narg('programme'), programme),
    au         = COALESCE(sqlc.narg('au'), au),
    updated_at =
        CASE
            WHEN (NOT (sqlc.narg('code')::TEXT IS NULL AND
                       sqlc.narg('year')::INTEGER IS NULL AND
                       sqlc.narg('semester')::TEXT IS NULL AND
                       sqlc.narg('programme')::TEXT IS NULL AND
                       sqlc.narg('au')::SMALLINT IS NULL))
                AND
                 (COALESCE(sqlc.narg('code'), code) <> code OR
                  COALESCE(sqlc.narg('year'), year) <> year OR
                  COALESCE(sqlc.narg('semester'), semester) <> semester OR
                  COALESCE(sqlc.narg('programme'), programme) <> programme OR
                  COALESCE(sqlc.narg('au'), au) <> au)
                THEN NOW()
            ELSE updated_at
            END
WHERE id = $1
RETURNING id, code, year, semester, programme, au, updated_at;

-- name: DeleteClass :one
DELETE
FROM classes
WHERE id = $1
RETURNING *;

-- name: UpsertClasses :batchone
-- Use this sparingly. When in doubt, use the atomic INSERT and UPDATE statements instead.
INSERT INTO classes (code, year, semester, programme, au, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
ON CONFLICT
    ON CONSTRAINT ux_code_year_semester
    DO UPDATE SET programme  = $4,
                  au         = $5,
                  updated_at =
                      CASE
                          WHEN $4 <> classes.programme OR
                               $5 <> classes.au
                              THEN NOW()
                          ELSE classes.updated_at
                          END
RETURNING *;