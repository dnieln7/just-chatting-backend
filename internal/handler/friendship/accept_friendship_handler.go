package friendship

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dnieln7/just-chatting/internal/database/db"
	"github.com/dnieln7/just-chatting/internal/helpers"
	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/google/uuid"
	"net/http"
)

type AcceptFriendshipBody struct {
	UserID   uuid.UUID `json:"user_id"`
	FriendID uuid.UUID `json:"friend_id"`
}

func AcceptFriendshipHandler(writer http.ResponseWriter, request *http.Request, resources *server.Resources) {
	decoder := json.NewDecoder(request.Body)
	body := AcceptFriendshipBody{}
	err := decoder.Decode(&body)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse JSON: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	err = AcceptBidirectionalFriendship(resources, request.Context(), body.UserID, body.FriendID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not accept friendship: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
	} else {
		helpers.ResponseOK(writer)
	}
}

func AcceptBidirectionalFriendship(resources *server.Resources, ctx context.Context, userID uuid.UUID, friendID uuid.UUID) error {
	transaction, err := resources.ConnectionDb.Begin()

	if err != nil {
		return err
	}

	defer transaction.Rollback()

	queriesTransaction := resources.PostgresDb.WithTx(transaction)

	err = queriesTransaction.UpdateFriendshipStatus(ctx, db.UpdateFriendshipStatusParams{
		Status:   db.FriendshipStatusAccepted,
		UserID:   userID,
		FriendID: friendID,
	})

	if err != nil {
		return err
	}

	err = queriesTransaction.UpdateFriendshipStatus(ctx, db.UpdateFriendshipStatusParams{
		Status:   db.FriendshipStatusAccepted,
		UserID:   friendID,
		FriendID: userID,
	})

	if err != nil {
		return err
	}

	return transaction.Commit()
}
