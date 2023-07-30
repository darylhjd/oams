// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: class_groups.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createClassGroup = `-- name: CreateClassGroup :one
INSERT INTO class_groups (class_id, name, class_type, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING id, class_id, name, class_type, created_at
`

type CreateClassGroupParams struct {
	ClassID   int64     `json:"class_id"`
	Name      string    `json:"name"`
	ClassType ClassType `json:"class_type"`
}

type CreateClassGroupRow struct {
	ID        int64            `json:"id"`
	ClassID   int64            `json:"class_id"`
	Name      string           `json:"name"`
	ClassType ClassType        `json:"class_type"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

func (q *Queries) CreateClassGroup(ctx context.Context, arg CreateClassGroupParams) (CreateClassGroupRow, error) {
	row := q.db.QueryRow(ctx, createClassGroup, arg.ClassID, arg.Name, arg.ClassType)
	var i CreateClassGroupRow
	err := row.Scan(
		&i.ID,
		&i.ClassID,
		&i.Name,
		&i.ClassType,
		&i.CreatedAt,
	)
	return i, err
}

const deleteClassGroup = `-- name: DeleteClassGroup :one
DELETE
FROM class_groups
WHERE id = $1
RETURNING id, class_id, name, class_type, created_at, updated_at
`

func (q *Queries) DeleteClassGroup(ctx context.Context, id int64) (ClassGroup, error) {
	row := q.db.QueryRow(ctx, deleteClassGroup, id)
	var i ClassGroup
	err := row.Scan(
		&i.ID,
		&i.ClassID,
		&i.Name,
		&i.ClassType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getClassGroup = `-- name: GetClassGroup :one
SELECT id, class_id, name, class_type, created_at, updated_at
FROM class_groups
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetClassGroup(ctx context.Context, id int64) (ClassGroup, error) {
	row := q.db.QueryRow(ctx, getClassGroup, id)
	var i ClassGroup
	err := row.Scan(
		&i.ID,
		&i.ClassID,
		&i.Name,
		&i.ClassType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getClassGroupsByIDs = `-- name: GetClassGroupsByIDs :many
SELECT id, class_id, name, class_type, created_at, updated_at
FROM class_groups
WHERE id = ANY ($1::BIGSERIAL[])
ORDER BY id
`

func (q *Queries) GetClassGroupsByIDs(ctx context.Context, ids []int64) ([]ClassGroup, error) {
	rows, err := q.db.Query(ctx, getClassGroupsByIDs, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ClassGroup
	for rows.Next() {
		var i ClassGroup
		if err := rows.Scan(
			&i.ID,
			&i.ClassID,
			&i.Name,
			&i.ClassType,
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

const listClassGroups = `-- name: ListClassGroups :many
SELECT id, class_id, name, class_type, created_at, updated_at
FROM class_groups
ORDER BY class_id, name
`

func (q *Queries) ListClassGroups(ctx context.Context) ([]ClassGroup, error) {
	rows, err := q.db.Query(ctx, listClassGroups)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ClassGroup
	for rows.Next() {
		var i ClassGroup
		if err := rows.Scan(
			&i.ID,
			&i.ClassID,
			&i.Name,
			&i.ClassType,
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

const updateClassGroup = `-- name: UpdateClassGroup :one
UPDATE class_groups
SET class_id   = COALESCE($2, class_id),
    name       = COALESCE($3, name),
    class_type = COALESCE($4, class_type),
    updated_at =
        CASE
            WHEN (NOT ($2::BIGSERIAL IS NULL AND
                       $3::TEXT IS NULL AND
                       $4::CLASS_TYPE IS NULL))
                AND
                 (COALESCE($2, class_id) <> class_id OR
                  COALESCE($3, name) <> name OR
                  COALESCE($4, class_type) <> class_type)
                THEN NOW()
            ELSE updated_at END
WHERE id = $1
RETURNING id, class_id, name, class_type, updated_at
`

type UpdateClassGroupParams struct {
	ID        int64         `json:"id"`
	ClassID   pgtype.Int8   `json:"class_id"`
	Name      pgtype.Text   `json:"name"`
	ClassType NullClassType `json:"class_type"`
}

type UpdateClassGroupRow struct {
	ID        int64            `json:"id"`
	ClassID   int64            `json:"class_id"`
	Name      string           `json:"name"`
	ClassType ClassType        `json:"class_type"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

func (q *Queries) UpdateClassGroup(ctx context.Context, arg UpdateClassGroupParams) (UpdateClassGroupRow, error) {
	row := q.db.QueryRow(ctx, updateClassGroup,
		arg.ID,
		arg.ClassID,
		arg.Name,
		arg.ClassType,
	)
	var i UpdateClassGroupRow
	err := row.Scan(
		&i.ID,
		&i.ClassID,
		&i.Name,
		&i.ClassType,
		&i.UpdatedAt,
	)
	return i, err
}
