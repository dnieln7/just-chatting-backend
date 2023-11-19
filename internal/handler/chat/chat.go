package chat

import (
	"time"

	"github.com/dnieln7/just-chatting/internal/database/db"
)

type Chats struct {
	Data []Chat `json:"data"`
}

type Chat struct {
	ID        string      `json:"id"`
	Me        Participant `json:"me"`
	Creator   Participant `json:"creator"`
	Friend    Participant `json:"friend"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type Participant struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func chatsWithMeAsCreatorToChat(me db.TbUser, chat db.GetChatsByCreatorIdRow) Chat {
	return Chat{
		ID: chat.ChatID.String(),
		Me: Participant{
			ID:       me.ID.String(),
			Email:    me.Email,
			Username: me.Username,
		},
		Creator: Participant{
			ID:       me.ID.String(),
			Email:    me.Email,
			Username: me.Username,
		},
		Friend: Participant{
			ID:       chat.FriendID.String(),
			Email:    chat.FriendEmail,
			Username: chat.FriendUsername,
		},
		CreatedAt: chat.ChatCreatedAt,
		UpdatedAt: chat.ChatUpdatedAt,
	}
}

func chatsWithMeAsFriendToChat(me db.TbUser, chat db.GetChatsByFriendIdRow) Chat {
	return Chat{
		ID: chat.ChatID.String(),
		Me: Participant{
			ID:       me.ID.String(),
			Email:    me.Email,
			Username: me.Username,
		},
		Creator: Participant{
			ID:       chat.CreatorID.String(),
			Email:    chat.CreatorEmail,
			Username: chat.CreatorUsername,
		},
		Friend: Participant{
			ID:       chat.CreatorID.String(),
			Email:    chat.CreatorEmail,
			Username: chat.CreatorUsername,
		},
		CreatedAt: chat.ChatCreatedAt,
		UpdatedAt: chat.ChatUpdatedAt,
	}
}
