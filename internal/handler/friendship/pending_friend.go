package friendship

import "github.com/dnieln7/just-chatting/internal/database/db"

type PendingFriend struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func dbGetPendingFriendsByUserIdRowToPendingFriend(row db.GetPendingFriendsByUserIdRow) PendingFriend {
	return PendingFriend{
		ID:       row.CreatorID.String(),
		Email:    row.Email,
		Username: row.Username,
	}
}
