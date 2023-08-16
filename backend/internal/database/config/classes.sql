-- name: UpdateClass :one
UPDATE classes
SET code       = COALESCE(sqlc.narg('code'), code),
    year       = COALESCE(sqlc.narg('year'), year),
    semester   = COALESCE(sqlc.narg('semester'), semester),
    programme  = COALESCE(sqlc.narg('programme'), programme),
    au         = COALESCE(sqlc.narg('au'), au)
WHERE id = $1
RETURNING id, code, year, semester, programme, au, updated_at;

-- name: UpsertClasses :batchone
-- Insert a class into the database. If the class already exists, then only update the programme and au.
INSERT INTO classes (code, year, semester, programme, au, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
ON CONFLICT
    ON CONSTRAINT ux_code_year_semester
    DO UPDATE SET programme  = $4,
                  au         = $5
RETURNING *;