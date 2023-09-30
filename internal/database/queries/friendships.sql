-- name: CreateFriendship :one
INSERT INTO tb_friendships (id, creator_id, user_id, friend_id, status, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: DeleteFriendship :exec
DELETE
FROM tb_friendships
WHERE user_id = $1
  AND friend_id = $2;

-- name: GetFriendshipByUserIdAndFriendId :many
SELECT *
FROM tb_friendships
WHERE user_id = $1
  AND friend_id = $2;

-- name: UpdateFriendshipStatus :exec
UPDATE tb_friendships
SET status = $1
WHERE user_id = $2
  AND friend_id = $3;

-- name: GetPendingFriendsByUserId :many
SELECT tb_friendships.creator_id, tb_users.email, tb_users.username
FROM tb_friendships
         JOIN tb_users ON tb_friendships.user_id = tb_users.id
WHERE creator_id != $1
  AND friend_id = $1
  AND status = 'pending';

-- name: GetFriendsByUserId :many
SELECT tb_users.id, tb_users.email, tb_users.username
FROM tb_friendships
         JOIN tb_users ON tb_friendships.friend_id = tb_users.id
WHERE user_id = $1
  AND status = 'accepted';
