// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: class_group_sessions.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createClassGroupSession = `-- name: CreateClassGroupSession :one
INSERT INTO class_group_sessions (class_group_id, start_time, end_time, venue, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id, class_group_id, start_time, end_time, venue, created_at
`

type CreateClassGroupSessionParams struct {
	ClassGroupID int64            `json:"class_group_id"`
	StartTime    pgtype.Timestamp `json:"start_time"`
	EndTime      pgtype.Timestamp `json:"end_time"`
	Venue        string           `json:"venue"`
}

type CreateClassGroupSessionRow struct {
	ID           int64            `json:"id"`
	ClassGroupID int64            `json:"class_group_id"`
	StartTime    pgtype.Timestamp `json:"start_time"`
	EndTime      pgtype.Timestamp `json:"end_time"`
	Venue        string           `json:"venue"`
	CreatedAt    pgtype.Timestamp `json:"created_at"`
}

func (q *Queries) CreateClassGroupSession(ctx context.Context, arg CreateClassGroupSessionParams) (CreateClassGroupSessionRow, error) {
	row := q.db.QueryRow(ctx, createClassGroupSession,
		arg.ClassGroupID,
		arg.StartTime,
		arg.EndTime,
		arg.Venue,
	)
	var i CreateClassGroupSessionRow
	err := row.Scan(
		&i.ID,
		&i.ClassGroupID,
		&i.StartTime,
		&i.EndTime,
		&i.Venue,
		&i.CreatedAt,
	)
	return i, err
}

const getClassGroupSession = `-- name: GetClassGroupSession :one
SELECT id, class_group_id, start_time, end_time, venue, created_at, updated_at
FROM class_group_sessions
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetClassGroupSession(ctx context.Context, id int64) (ClassGroupSession, error) {
	row := q.db.QueryRow(ctx, getClassGroupSession, id)
	var i ClassGroupSession
	err := row.Scan(
		&i.ID,
		&i.ClassGroupID,
		&i.StartTime,
		&i.EndTime,
		&i.Venue,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getClassGroupSessionsByIDs = `-- name: GetClassGroupSessionsByIDs :many
SELECT id, class_group_id, start_time, end_time, venue, created_at, updated_at
FROM class_group_sessions
WHERE id = ANY ($1::BIGINT[])
ORDER BY id
`

func (q *Queries) GetClassGroupSessionsByIDs(ctx context.Context, ids []int64) ([]ClassGroupSession, error) {
	rows, err := q.db.Query(ctx, getClassGroupSessionsByIDs, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ClassGroupSession
	for rows.Next() {
		var i ClassGroupSession
		if err := rows.Scan(
			&i.ID,
			&i.ClassGroupID,
			&i.StartTime,
			&i.EndTime,
			&i.Venue,
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

const listClassGroupSessions = `-- name: ListClassGroupSessions :many
SELECT id, class_group_id, start_time, end_time, venue, created_at, updated_at
FROM class_group_sessions
ORDER BY class_group_id, start_time
`

func (q *Queries) ListClassGroupSessions(ctx context.Context) ([]ClassGroupSession, error) {
	rows, err := q.db.Query(ctx, listClassGroupSessions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ClassGroupSession
	for rows.Next() {
		var i ClassGroupSession
		if err := rows.Scan(
			&i.ID,
			&i.ClassGroupID,
			&i.StartTime,
			&i.EndTime,
			&i.Venue,
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

const updateClassGroupSession = `-- name: UpdateClassGroupSession :one
UPDATE class_group_sessions
SET class_group_id = COALESCE($2, class_group_id),
    start_time     = COALESCE($3, start_time),
    end_time       = COALESCE($4, end_time),
    venue          = COALESCE($5, venue),
    updated_at     =
        CASE
            WHEN (NOT ($2::BIGINT IS NULL AND
                       $3::TIMESTAMP IS NULL AND
                       $4::TIMESTAMP IS NULL AND
                       $5::TEXT IS NULL))
                AND
                 (COALESCE($2, class_group_id) <> class_group_id OR
                  COALESCE($3, start_time) <> start_time OR
                  COALESCE($4, end_time) <> end_time OR
                  COALESCE($5, venue) <> venue)
                THEN NOW()
            ELSE updated_at
            END
WHERE id = $1
RETURNING id, class_group_id, start_time, end_time, venue, updated_at
`

type UpdateClassGroupSessionParams struct {
	ID           int64            `json:"id"`
	ClassGroupID pgtype.Int8      `json:"class_group_id"`
	StartTime    pgtype.Timestamp `json:"start_time"`
	EndTime      pgtype.Timestamp `json:"end_time"`
	Venue        pgtype.Text      `json:"venue"`
}

type UpdateClassGroupSessionRow struct {
	ID           int64            `json:"id"`
	ClassGroupID int64            `json:"class_group_id"`
	StartTime    pgtype.Timestamp `json:"start_time"`
	EndTime      pgtype.Timestamp `json:"end_time"`
	Venue        string           `json:"venue"`
	UpdatedAt    pgtype.Timestamp `json:"updated_at"`
}

func (q *Queries) UpdateClassGroupSession(ctx context.Context, arg UpdateClassGroupSessionParams) (UpdateClassGroupSessionRow, error) {
	row := q.db.QueryRow(ctx, updateClassGroupSession,
		arg.ID,
		arg.ClassGroupID,
		arg.StartTime,
		arg.EndTime,
		arg.Venue,
	)
	var i UpdateClassGroupSessionRow
	err := row.Scan(
		&i.ID,
		&i.ClassGroupID,
		&i.StartTime,
		&i.EndTime,
		&i.Venue,
		&i.UpdatedAt,
	)
	return i, err
}
