// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: query.sql

package database

import (
	"context"
)

const getStudent = `-- name: GetStudent :one
SELECT matric_no, name, email
FROM Students
WHERE matric_no = $1
LIMIT 1
`

func (q *Queries) GetStudent(ctx context.Context, matricNo string) (Student, error) {
	row := q.db.QueryRowContext(ctx, getStudent, matricNo)
	var i Student
	err := row.Scan(&i.MatricNo, &i.Name, &i.Email)
	return i, err
}

const listStudents = `-- name: ListStudents :many
SELECT matric_no, name, email
FROM Students
ORDER BY matric_no
`

func (q *Queries) ListStudents(ctx context.Context) ([]Student, error) {
	rows, err := q.db.QueryContext(ctx, listStudents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Student
	for rows.Next() {
		var i Student
		if err := rows.Scan(&i.MatricNo, &i.Name, &i.Email); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
