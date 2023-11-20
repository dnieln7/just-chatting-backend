package chat

import (
	"context"
	"fmt"
	"github.com/dnieln7/just-chatting/internal/database/db"
	"github.com/google/uuid"
)

func getCreatorAndFriend(
	postgresDb *db.Queries,
	context context.Context,
	creatorID uuid.UUID,
	friendID uuid.UUID,
) (errMessage string, creator db.TbUser, friend db.TbUser) {
	c, err := postgresDb.GetUserById(context, creatorID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not find user: %v", err)
		return errMessage, db.TbUser{}, db.TbUser{}
	}

	f, err := postgresDb.GetUserById(context, friendID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not find user: %v", err)
		return errMessage, db.TbUser{}, db.TbUser{}
	}

	return "", c, f
}
