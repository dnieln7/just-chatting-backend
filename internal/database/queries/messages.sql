-- name: CreateMessage :one
INSERT INTO tb_messages (
        id,
        chat_id,
        user_id,
        message,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetMessagesByChatId :many
SELECT * FROM tb_messages WHERE chat_id = $1 ORDER BY created_at DESC;

-- name: GetMessagesByChatIdLazy :many
SELECT * FROM tb_messages WHERE chat_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3;
