package friendship

import (
	"fmt"
	"github.com/dnieln7/just-chatting/internal/helpers"
	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

func GetFriendsHandler(writer http.ResponseWriter, request *http.Request, resources *server.Resources) {
	vars := mux.Vars(request)
	stringUUID := vars["id"]

	userID, err := uuid.Parse(stringUUID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse UUID: %v", stringUUID)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	rows, err := resources.PostgresDb.GetFriendsByUserId(request.Context(), userID)

	if err != nil {
		errMessage := fmt.Sprintf("Could get friends: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
	} else {
		friends := []Friend{}

		for _, row := range rows {
			friends = append(friends, dbGetFriendsByUserIdRowToFriend(row))
		}

		helpers.ResponseJson(writer, 200, friends)
	}
}
