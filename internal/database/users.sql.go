// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO
  users (id, created_at, updated_at, name, apikey)
VALUES
  (
    $1,
    $2,
    $3,
    $4,
    encode(
      sha256(random()::text::bytea),
      'hex'
    )
  ) RETURNING id, created_at, updated_at, name, apikey
`

type CreateUserParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Apikey,
	)
	return i, err
}

const getUserByApiKey = `-- name: GetUserByApiKey :one
SELECT
  id, created_at, updated_at, name, apikey
FROM
  users
WHERE
  apikey = $1
`

func (q *Queries) GetUserByApiKey(ctx context.Context, apikey string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByApiKey, apikey)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Apikey,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT
  id, created_at, updated_at, name, apikey
FROM
  users
`

func (q *Queries) GetUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Apikey,
		); err != nil {
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

const removeAllUsers = `-- name: RemoveAllUsers :exec
DELETE FROM
  users
`

func (q *Queries) RemoveAllUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, removeAllUsers)
	return err
}
