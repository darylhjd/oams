-- name: GetStudent :one
SELECT *
FROM Students
WHERE matric_no = $1
LIMIT 1;

-- name: ListStudents :many
SELECT *
FROM Students
ORDER BY matric_no;