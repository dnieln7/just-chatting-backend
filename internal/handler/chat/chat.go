package chat

import (
	"time"

	"github.com/dnieln7/just-chatting/internal/database/db"
	"github.com/dnieln7/just-chatting/internal/helpers"
)

type Chats struct {
	Data []Chat `json:"data"`
}

type Chat struct {
	ID           string    `json:"id"`
	Participants []string  `json:"participants"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func dbChatToChat(dbChat db.TbChat) Chat {
	return Chat{
		ID:           dbChat.ID.String(),
		Participants: helpers.UUIDsToStrings(dbChat.Participants),
		CreatedAt:    dbChat.CreatedAt,
		UpdatedAt:    dbChat.UpdatedAt,
	}
}
