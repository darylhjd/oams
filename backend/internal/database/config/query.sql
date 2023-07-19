-- name: GetStudent :one
SELECT *
FROM Students
WHERE id = $1
LIMIT 1;

-- name: ListStudents :many
SELECT *
FROM Students
ORDER BY id;