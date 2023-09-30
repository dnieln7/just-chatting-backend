-- +goose Up
CREATE TYPE friendship_status AS ENUM ('pending','accepted');

CREATE TABLE tb_friendships
(
    id         UUID PRIMARY KEY,
    creator_id UUID              NOT NULL,
    user_id    UUID              NOT NULL REFERENCES tb_users (id) ON DELETE CASCADE,
    friend_id  UUID              NOT NULL REFERENCES tb_users (id) ON DELETE CASCADE,
    status     friendship_status NOT NULL,
    created_at TIMESTAMP         NOT NULL,
    updated_at TIMESTAMP         NOT NULL,
    UNIQUE (user_id, friend_id)
);

-- +goose Down
DROP TABLE tb_friendships;

DROP TYPE friendship_status;
