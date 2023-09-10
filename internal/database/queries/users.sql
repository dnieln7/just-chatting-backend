-- name: CreateUser :one
INSERT INTO tb_users (
        id,
        email,
        password,
        username,
        created_at,
        updated_at
    )
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM tb_users WHERE email = $1;