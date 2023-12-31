// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: friendships.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createFriendship = `-- name: CreateFriendship :one
INSERT INTO tb_friendships (id, creator_id, user_id, friend_id, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, creator_id, user_id, friend_id, status, created_at, updated_at
`

type CreateFriendshipParams struct {
	ID        uuid.UUID
	CreatorID uuid.UUID
	UserID    uuid.UUID
	FriendID  uuid.UUID
	Status    FriendshipStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateFriendship(ctx context.Context, arg CreateFriendshipParams) (TbFriendship, error) {
	row := q.db.QueryRowContext(ctx, createFriendship,
		arg.ID,
		arg.CreatorID,
		arg.UserID,
		arg.FriendID,
		arg.Status,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i TbFriendship
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.UserID,
		&i.FriendID,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteFriendship = `-- name: DeleteFriendship :exec
DELETE
FROM tb_friendships
WHERE user_id = $1
  AND friend_id = $2
`

type DeleteFriendshipParams struct {
	UserID   uuid.UUID
	FriendID uuid.UUID
}

func (q *Queries) DeleteFriendship(ctx context.Context, arg DeleteFriendshipParams) error {
	_, err := q.db.ExecContext(ctx, deleteFriendship, arg.UserID, arg.FriendID)
	return err
}

const getFriendsByUserId = `-- name: GetFriendsByUserId :many
SELECT tb_users.id, tb_users.email, tb_users.username
FROM tb_friendships
         JOIN tb_users ON tb_friendships.friend_id = tb_users.id
WHERE user_id = $1
  AND status = 'accepted'
`

type GetFriendsByUserIdRow struct {
	ID       uuid.UUID
	Email    string
	Username string
}

func (q *Queries) GetFriendsByUserId(ctx context.Context, userID uuid.UUID) ([]GetFriendsByUserIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getFriendsByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFriendsByUserIdRow
	for rows.Next() {
		var i GetFriendsByUserIdRow
		if err := rows.Scan(&i.ID, &i.Email, &i.Username); err != nil {
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

const getFriendshipByUserIdAndFriendId = `-- name: GetFriendshipByUserIdAndFriendId :many
SELECT id, creator_id, user_id, friend_id, status, created_at, updated_at
FROM tb_friendships
WHERE user_id = $1
  AND friend_id = $2
`

type GetFriendshipByUserIdAndFriendIdParams struct {
	UserID   uuid.UUID
	FriendID uuid.UUID
}

func (q *Queries) GetFriendshipByUserIdAndFriendId(ctx context.Context, arg GetFriendshipByUserIdAndFriendIdParams) ([]TbFriendship, error) {
	rows, err := q.db.QueryContext(ctx, getFriendshipByUserIdAndFriendId, arg.UserID, arg.FriendID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TbFriendship
	for rows.Next() {
		var i TbFriendship
		if err := rows.Scan(
			&i.ID,
			&i.CreatorID,
			&i.UserID,
			&i.FriendID,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const getPendingFriendsByUserId = `-- name: GetPendingFriendsByUserId :many
SELECT tb_friendships.creator_id, tb_users.email, tb_users.username
FROM tb_friendships
         JOIN tb_users ON tb_friendships.user_id = tb_users.id
WHERE creator_id != $1
  AND friend_id = $1
  AND status = 'pending'
`

type GetPendingFriendsByUserIdRow struct {
	CreatorID uuid.UUID
	Email     string
	Username  string
}

func (q *Queries) GetPendingFriendsByUserId(ctx context.Context, creatorID uuid.UUID) ([]GetPendingFriendsByUserIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getPendingFriendsByUserId, creatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPendingFriendsByUserIdRow
	for rows.Next() {
		var i GetPendingFriendsByUserIdRow
		if err := rows.Scan(&i.CreatorID, &i.Email, &i.Username); err != nil {
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

const updateFriendshipStatus = `-- name: UpdateFriendshipStatus :exec
UPDATE tb_friendships
SET status = $1
WHERE user_id = $2
  AND friend_id = $3
`

type UpdateFriendshipStatusParams struct {
	Status   FriendshipStatus
	UserID   uuid.UUID
	FriendID uuid.UUID
}

func (q *Queries) UpdateFriendshipStatus(ctx context.Context, arg UpdateFriendshipStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateFriendshipStatus, arg.Status, arg.UserID, arg.FriendID)
	return err
}
