package user

import (
	"time"

	"github.com/dnieln7/just-chatting/internal/database/db"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func dbUserToUser(dbUser db.TbUser) User {
	return User{
		ID:        dbUser.ID.String(),
		Email:     dbUser.Email,
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}
}
