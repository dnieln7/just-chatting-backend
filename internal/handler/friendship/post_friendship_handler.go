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
	"time"
)

type PostFriendshipBody struct {
	UserID   uuid.UUID `json:"user_id"`
	FriendID uuid.UUID `json:"friend_id"`
}

func PostFriendshipHandler(writer http.ResponseWriter, request *http.Request, resources *server.Resources) {
	decoder := json.NewDecoder(request.Body)
	body := PostFriendshipBody{}
	err := decoder.Decode(&body)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse JSON: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	dbFriendship, err := resources.PostgresDb.GetFriendshipByUserIdAndFriendId(request.Context(), db.GetFriendshipByUserIdAndFriendIdParams{
		UserID:   body.UserID,
		FriendID: body.FriendID,
	})

	if len(dbFriendship) > 0 {
		helpers.ResponseJsonError(writer, 409, "Friendship is already created with status pending or accepted")
		return
	}

	err = CreateBidirectionalFriendship(resources, request.Context(), body.UserID, body.FriendID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not create friendship: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
	} else {
		helpers.ResponseNoContent(writer)
	}
}

func CreateBidirectionalFriendship(resources *server.Resources, ctx context.Context, userID uuid.UUID, friendID uuid.UUID) error {
	transaction, err := resources.ConnectionDb.Begin()

	if err != nil {
		return err
	}

	defer transaction.Rollback()

	queriesTransaction := resources.PostgresDb.WithTx(transaction)

	_, err = queriesTransaction.CreateFriendship(ctx, db.CreateFriendshipParams{
		ID:        uuid.New(),
		CreatorID: userID,
		UserID:    userID,
		FriendID:  friendID,
		Status:    db.FriendshipStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return err
	}

	_, err = queriesTransaction.CreateFriendship(ctx, db.CreateFriendshipParams{
		ID:        uuid.New(),
		CreatorID: userID,
		UserID:    friendID,
		FriendID:  userID,
		Status:    db.FriendshipStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		return err
	}

	return transaction.Commit()
}
