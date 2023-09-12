-- +goose Up
CREATE TABLE tb_chats (
    id UUID PRIMARY KEY,
    participants UUID[] NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE tb_chats;
