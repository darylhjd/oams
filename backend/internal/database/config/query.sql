-- name: ListStudents :many
SELECT *
FROM students
ORDER BY id;

-- name: GetStudent :one
SELECT *
FROM students
WHERE id = $1
LIMIT 1;

-- name: UpsertStudents :batchone
INSERT INTO students (id, name, email, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT (id)
    DO UPDATE SET name       = $2,
                  email      = $3,
                  updated_at = NOW()
RETURNING *;

-- name: ListCourses :many
SELECT *
FROM courses
ORDER BY code, year, semester;

-- name: UpsertCourses :batchone
INSERT INTO courses (code, year, semester, programme, au, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
ON CONFLICT ON CONSTRAINT ux_code_year_semester
    DO UPDATE SET programme  = $4,
                  au         = $5,
                  updated_at = NOW()
RETURNING *;

-- name: ListClassGroups :many
SELECT *
FROM class_groups
ORDER BY course_id, name;

-- name: UpsertClassGroups :batchone
INSERT INTO class_groups (course_id, name, class_type, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT ON CONSTRAINT ux_course_id_name
    DO UPDATE SET class_type = $3,
                  updated_at = NOW()
RETURNING *;

-- name: ListClassGroupSessions :many
SELECT *
FROM class_group_sessions
ORDER BY class_group_id, start_time, end_time;

-- name: UpsertClassGroupSessions :batchone
INSERT INTO class_group_sessions (class_group_id, start_time, end_time, venue, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
ON CONFLICT ON CONSTRAINT ux_class_group_id_start_time
    DO UPDATE SET end_time   = $3,
                  venue      = $4,
                  updated_at = NOW()
RETURNING *;

-- name: ListSessionEnrollments :many
SELECT *
FROM session_enrollments
ORDER BY session_id, student_id;

-- name: CreateSessionEnrollments :batchone
INSERT INTO session_enrollments (session_id, student_id, created_at)
VALUES ($1, $2, NOW())
RETURNING *;