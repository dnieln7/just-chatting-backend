// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: messages.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createMessage = `-- name: CreateMessage :one
INSERT INTO tb_messages (
        id,
        chat_id,
        user_id,
        message,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, chat_id, user_id, message, created_at, updated_at
`

type CreateMessageParams struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	UserID    uuid.UUID
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) (TbMessage, error) {
	row := q.db.QueryRowContext(ctx, createMessage,
		arg.ID,
		arg.ChatID,
		arg.UserID,
		arg.Message,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i TbMessage
	err := row.Scan(
		&i.ID,
		&i.ChatID,
		&i.UserID,
		&i.Message,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getMessagesByChatId = `-- name: GetMessagesByChatId :many
SELECT id, chat_id, user_id, message, created_at, updated_at FROM tb_messages WHERE chat_id = $1 ORDER BY created_at DESC
`

func (q *Queries) GetMessagesByChatId(ctx context.Context, chatID uuid.UUID) ([]TbMessage, error) {
	rows, err := q.db.QueryContext(ctx, getMessagesByChatId, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TbMessage
	for rows.Next() {
		var i TbMessage
		if err := rows.Scan(
			&i.ID,
			&i.ChatID,
			&i.UserID,
			&i.Message,
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

const getMessagesByChatIdLazy = `-- name: GetMessagesByChatIdLazy :many
SELECT id, chat_id, user_id, message, created_at, updated_at FROM tb_messages WHERE chat_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
`

type GetMessagesByChatIdLazyParams struct {
	ChatID uuid.UUID
	Limit  int32
	Offset int32
}

func (q *Queries) GetMessagesByChatIdLazy(ctx context.Context, arg GetMessagesByChatIdLazyParams) ([]TbMessage, error) {
	rows, err := q.db.QueryContext(ctx, getMessagesByChatIdLazy, arg.ChatID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TbMessage
	for rows.Next() {
		var i TbMessage
		if err := rows.Scan(
			&i.ID,
			&i.ChatID,
			&i.UserID,
			&i.Message,
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
