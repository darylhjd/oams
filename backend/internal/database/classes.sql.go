// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: classes.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createClass = `-- name: CreateClass :one
INSERT INTO classes (code, year, semester, programme, au, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
RETURNING id, code, year, semester, programme, au, created_at
`

type CreateClassParams struct {
	Code      string `json:"code"`
	Year      int32  `json:"year"`
	Semester  string `json:"semester"`
	Programme string `json:"programme"`
	Au        int16  `json:"au"`
}

type CreateClassRow struct {
	ID        int64            `json:"id"`
	Code      string           `json:"code"`
	Year      int32            `json:"year"`
	Semester  string           `json:"semester"`
	Programme string           `json:"programme"`
	Au        int16            `json:"au"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

func (q *Queries) CreateClass(ctx context.Context, arg CreateClassParams) (CreateClassRow, error) {
	row := q.db.QueryRow(ctx, createClass,
		arg.Code,
		arg.Year,
		arg.Semester,
		arg.Programme,
		arg.Au,
	)
	var i CreateClassRow
	err := row.Scan(
		&i.ID,
		&i.Code,
		&i.Year,
		&i.Semester,
		&i.Programme,
		&i.Au,
		&i.CreatedAt,
	)
	return i, err
}

const deleteClass = `-- name: DeleteClass :one
DELETE
FROM classes
WHERE id = $1
RETURNING id, code, year, semester, programme, au, created_at, updated_at
`

func (q *Queries) DeleteClass(ctx context.Context, id int64) (Class, error) {
	row := q.db.QueryRow(ctx, deleteClass, id)
	var i Class
	err := row.Scan(
		&i.ID,
		&i.Code,
		&i.Year,
		&i.Semester,
		&i.Programme,
		&i.Au,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getClass = `-- name: GetClass :one
SELECT id, code, year, semester, programme, au, created_at, updated_at
FROM classes
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetClass(ctx context.Context, id int64) (Class, error) {
	row := q.db.QueryRow(ctx, getClass, id)
	var i Class
	err := row.Scan(
		&i.ID,
		&i.Code,
		&i.Year,
		&i.Semester,
		&i.Programme,
		&i.Au,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getClassesByIDs = `-- name: GetClassesByIDs :many
SELECT id, code, year, semester, programme, au, created_at, updated_at
FROM classes
WHERE id = ANY ($1::BIGINT[])
ORDER BY id
`

func (q *Queries) GetClassesByIDs(ctx context.Context, ids []int64) ([]Class, error) {
	rows, err := q.db.Query(ctx, getClassesByIDs, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Class
	for rows.Next() {
		var i Class
		if err := rows.Scan(
			&i.ID,
			&i.Code,
			&i.Year,
			&i.Semester,
			&i.Programme,
			&i.Au,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listClasses = `-- name: ListClasses :many
SELECT id, code, year, semester, programme, au, created_at, updated_at
FROM classes
ORDER BY code, year, semester
`

func (q *Queries) ListClasses(ctx context.Context) ([]Class, error) {
	rows, err := q.db.Query(ctx, listClasses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Class
	for rows.Next() {
		var i Class
		if err := rows.Scan(
			&i.ID,
			&i.Code,
			&i.Year,
			&i.Semester,
			&i.Programme,
			&i.Au,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateClass = `-- name: UpdateClass :one
UPDATE classes
SET code       = COALESCE($2, code),
    year       = COALESCE($3, year),
    semester   = COALESCE($4, semester),
    programme  = COALESCE($5, programme),
    au         = COALESCE($6, au),
    updated_at =
        CASE
            WHEN (NOT ($2::TEXT IS NULL AND
                       $3::INTEGER IS NULL AND
                       $4::TEXT IS NULL AND
                       $5::TEXT IS NULL AND
                       $6::SMALLINT IS NULL))
                AND
                 (COALESCE($2, code) <> code OR
                  COALESCE($3, year) <> year OR
                  COALESCE($4, semester) <> semester OR
                  COALESCE($5, programme) <> programme OR
                  COALESCE($6, au) <> au)
                THEN NOW()
            ELSE updated_at
            END
WHERE id = $1
RETURNING id, code, year, semester, programme, au, updated_at
`

type UpdateClassParams struct {
	ID        int64       `json:"id"`
	Code      pgtype.Text `json:"code"`
	Year      pgtype.Int4 `json:"year"`
	Semester  pgtype.Text `json:"semester"`
	Programme pgtype.Text `json:"programme"`
	Au        pgtype.Int2 `json:"au"`
}

type UpdateClassRow struct {
	ID        int64            `json:"id"`
	Code      string           `json:"code"`
	Year      int32            `json:"year"`
	Semester  string           `json:"semester"`
	Programme string           `json:"programme"`
	Au        int16            `json:"au"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

func (q *Queries) UpdateClass(ctx context.Context, arg UpdateClassParams) (UpdateClassRow, error) {
	row := q.db.QueryRow(ctx, updateClass,
		arg.ID,
		arg.Code,
		arg.Year,
		arg.Semester,
		arg.Programme,
		arg.Au,
	)
	var i UpdateClassRow
	err := row.Scan(
		&i.ID,
		&i.Code,
		&i.Year,
		&i.Semester,
		&i.Programme,
		&i.Au,
		&i.UpdatedAt,
	)
	return i, err
}
