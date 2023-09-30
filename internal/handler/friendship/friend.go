package friendship

import "github.com/dnieln7/just-chatting/internal/database/db"

type Friend struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func dbGetFriendsByUserIdRowToFriend(row db.GetFriendsByUserIdRow) Friend {
	return Friend{
		ID:       row.ID.String(),
		Email:    row.Email,
		Username: row.Username,
	}
}
