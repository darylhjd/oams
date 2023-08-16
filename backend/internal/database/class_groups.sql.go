// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: class_groups.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const updateClassGroup = `-- name: UpdateClassGroup :one
UPDATE class_groups
SET class_id   = COALESCE($2, class_id),
    name       = COALESCE($3, name),
    class_type = COALESCE($4, class_type),
    updated_at =
        CASE
            WHEN (COALESCE($2, class_id) <> class_id OR
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
	ID        int64              `json:"id"`
	ClassID   int64              `json:"class_id"`
	Name      string             `json:"name"`
	ClassType ClassType          `json:"class_type"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
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
