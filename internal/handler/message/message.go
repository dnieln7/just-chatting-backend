package message

import (
	"time"

	"github.com/dnieln7/just-chatting/internal/database/db"
)

type Messages struct {
	Data        []Message `json:"data"`
	CurrentPage int       `json:"current_page"`
	HasNextPage bool      `json:"has_next_page"`
}

type Message struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chat_id"`
	UserID    string    `json:"user_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DBMessageToMessage(dbMessage db.TbMessage) Message {
	return Message{
		ID:        dbMessage.ID.String(),
		ChatID:    dbMessage.ChatID.String(),
		UserID:    dbMessage.UserID.String(),
		Message:   dbMessage.Message,
		CreatedAt: dbMessage.CreatedAt,
		UpdatedAt: dbMessage.UpdatedAt,
	}
}
