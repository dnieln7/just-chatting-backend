// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: users.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const countUsersWithEmail = `-- name: CountUsersWithEmail :one
SELECT count(*)
FROM tb_users
WHERE email = $1
`

func (q *Queries) CountUsersWithEmail(ctx context.Context, email string) (int64, error) {
	row := q.db.QueryRowContext(ctx, countUsersWithEmail, email)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO tb_users (
        id,
        email,
        password,
        username,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, email, password, username, created_at, updated_at
`

type CreateUserParams struct {
	ID        uuid.UUID
	Email     string
	Password  string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (TbUser, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.Password,
		arg.Username,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i TbUser
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.Username,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, password, username, created_at, updated_at FROM tb_users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (TbUser, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i TbUser
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.Username,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, email, password, username, created_at, updated_at
FROM tb_users
WHERE id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (TbUser, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i TbUser
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.Username,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
