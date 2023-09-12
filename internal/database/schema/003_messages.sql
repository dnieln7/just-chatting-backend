-- +goose Up
CREATE TABLE tb_messages (
    id UUID PRIMARY KEY,
    chat_id UUID NOT NULL REFERENCES tb_chats(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES tb_users(id),
    message TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE tb_messages;