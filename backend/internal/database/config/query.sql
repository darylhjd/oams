-- name: ListStudents :many
SELECT *
FROM students
ORDER BY id;

-- name: GetStudent :one
SELECT *
FROM students
WHERE id = $1
LIMIT 1;

-- name: UpsertStudents :batchmany
INSERT INTO students (id, name, email, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT (id)
DO UPDATE SET name = $2, email = $3, updated_at = NOW()
RETURNING id;

-- name: CreateCourses :batchmany
INSERT INTO courses (code, year, semester, programme, au, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
RETURNING id;

-- name: CreateClassGroups :batchmany
INSERT INTO class_groups (course_id, name, class_type, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING id;

-- name: CreateClassGroupSessions :batchmany
INSERT INTO class_group_sessions (class_group_id, start_time, end_time, venue, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id;

-- name: CreateSessionEnrollments :batchmany
INSERT INTO session_enrollments (session_id, student_id)
VALUES ($1, $2)
RETURNING session_id, student_id;