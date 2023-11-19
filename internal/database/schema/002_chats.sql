-- +goose Up
CREATE TABLE tb_chats
(
    id         UUID PRIMARY KEY,
    creator_id UUID      NOT NULL REFERENCES tb_users (id) ON DELETE CASCADE,
    friend_id  UUID      NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE (creator_id, friend_id)
);

-- +goose Down
DROP TABLE tb_chats;
