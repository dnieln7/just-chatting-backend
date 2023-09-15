-- name: CreateChat :one
INSERT INTO tb_chats (
        id,
        participants,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetChatsByParticipantId :many
SELECT * FROM tb_chats WHERE participants @> $1;

-- name: GetChatWithParticipants :one
SELECT * FROM tb_chats WHERE participants @> $1 LIMIT 1;

-- name: GetChatById :one
SELECT * FROM tb_chats WHERE id = $1;

-- name: GetParticipantsByChatId :one
SELECT participants FROM tb_chats WHERE id = $1;
