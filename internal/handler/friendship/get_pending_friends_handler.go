package friendship

import (
	"fmt"
	"github.com/dnieln7/just-chatting/internal/helpers"
	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

func GetPendingFriendsHandler(writer http.ResponseWriter, request *http.Request, resources *server.Resources) {
	vars := mux.Vars(request)
	stringUUID := vars["id"]

	userID, err := uuid.Parse(stringUUID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse UUID: %v", stringUUID)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	rows, err := resources.PostgresDb.GetPendingFriendsByUserId(request.Context(), userID)

	if err != nil {
		errMessage := fmt.Sprintf("Could get pending friends: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
	} else {
		pendingFriends := []PendingFriend{}

		for _, row := range rows {
			pendingFriends = append(pendingFriends, dbGetPendingFriendsByUserIdRowToPendingFriend(row))
		}

		helpers.ResponseJson(writer, 200, pendingFriends)
	}
}
