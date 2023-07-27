// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: users.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, name, email, role, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id, name, email, role, created_at
`

type CreateUserParams struct {
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	Email pgtype.Text `json:"email"`
	Role  UserRole    `json:"role"`
}

type CreateUserRow struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Email     pgtype.Text      `json:"email"`
	Role      UserRole         `json:"role"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Role,
	)
	var i CreateUserRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Role,
		&i.CreatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, name, email, role, created_at, updated_at
FROM users
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUsersByIDs = `-- name: GetUsersByIDs :many
SELECT id, name, email, role, created_at, updated_at
FROM users
WHERE id = ANY ($1::TEXT[])
ORDER BY id
`

func (q *Queries) GetUsersByIDs(ctx context.Context, ids []string) ([]User, error) {
	rows, err := q.db.Query(ctx, getUsersByIDs, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Role,
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

const listUsers = `-- name: ListUsers :many
SELECT id, name, email, role, created_at, updated_at
FROM users
ORDER BY id
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Email,
			&i.Role,
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

const updateUser = `-- name: UpdateUser :one
UPDATE users
    SET name = COALESCE($2, name),
        email = COALESCE($3, email),
        role = COALESCE($4, role),
        updated_at = NOW()
WHERE id = $1
RETURNING id, name, email, role, created_at, updated_at
`

type UpdateUserParams struct {
	ID    string       `json:"id"`
	Name  pgtype.Text  `json:"name"`
	Email pgtype.Text  `json:"email"`
	Role  NullUserRole `json:"role"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.ID,
		arg.Name,
		arg.Email,
		arg.Role,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
