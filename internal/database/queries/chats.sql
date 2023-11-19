-- name: CreateChat :one
INSERT INTO tb_chats (id,
                      creator_id,
                      friend_id,
                      created_at,
                      updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetChatById :one
SELECT *
FROM tb_chats
WHERE id = $1;

-- name: GetChatsByUserId :many
SELECT *
FROM tb_chats
WHERE creator_id = $1
   OR friend_id = $1;

-- name: GetChatWithCreatorAndFriend :one
SELECT *
FROM tb_chats
WHERE (creator_id = $1 AND friend_id = $2)
   OR (creator_id = $2 AND friend_id = $1)
LIMIT 1;

-- name: GetChatsByCreatorId :many
SELECT tb_chats.id         AS chat_id,
       tb_chats.created_at AS chat_created_at,
       tb_chats.updated_at AS chat_updated_at,
       tb_chats.creator_id AS chat_creator_id,
       tb_users.id         AS friend_id,
       tb_users.email      AS friend_email,
       tb_users.username   AS friend_username
FROM tb_chats
         JOIN tb_users ON tb_users.id = tb_chats.friend_id
WHERE tb_chats.creator_id = $1;

-- name: GetChatsByFriendId :many
SELECT tb_chats.id         AS chat_id,
       tb_chats.created_at AS chat_created_at,
       tb_chats.updated_at AS chat_updated_at,
       tb_chats.friend_id  AS chat_friend_id,
       tb_users.id         AS creator_id,
       tb_users.email      AS creator_email,
       tb_users.username   AS creator_username
FROM tb_chats
         JOIN tb_users ON tb_users.id = tb_chats.creator_id
WHERE tb_chats.friend_id = $1;
