package friendship

import (
	"context"
	"fmt"
	"github.com/dnieln7/just-chatting/internal/database/db"
	"github.com/dnieln7/just-chatting/internal/helpers"
	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

func DeleteFriendshipHandler(writer http.ResponseWriter, request *http.Request, resources *server.Resources) {
	vars := mux.Vars(request)
	stringUserUUID := vars["user_id"]
	stringFriendUUID := vars["friend_id"]

	userID, err := uuid.Parse(stringUserUUID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse UUID: %v", stringUserUUID)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	friendID, err := uuid.Parse(stringFriendUUID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse UUID: %v", stringFriendUUID)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	err = DeleteBidirectionalFriendship(resources, request.Context(), userID, friendID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not delete friendship: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
	} else {
		helpers.ResponseOK(writer)
	}
}

func DeleteBidirectionalFriendship(resources *server.Resources, ctx context.Context, userID uuid.UUID, friendID uuid.UUID) error {
	transaction, err := resources.ConnectionDb.Begin()

	if err != nil {
		return err
	}

	defer transaction.Rollback()

	queriesTransaction := resources.PostgresDb.WithTx(transaction)

	err = queriesTransaction.DeleteFriendship(ctx, db.DeleteFriendshipParams{
		UserID:   userID,
		FriendID: friendID,
	})

	if err != nil {
		return err
	}

	err = queriesTransaction.DeleteFriendship(ctx, db.DeleteFriendshipParams{
		UserID:   friendID,
		FriendID: userID,
	})

	if err != nil {
		return err
	}

	return transaction.Commit()
}
