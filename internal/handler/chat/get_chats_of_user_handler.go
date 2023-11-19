package chat

import (
	"fmt"
	"net/http"

	"github.com/dnieln7/just-chatting/internal/helpers"
	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetChatsOfUserHandler(writer http.ResponseWriter, request *http.Request, resources *server.Resources) {
	vars := mux.Vars(request)
	stringUUID := vars["id"]

	userID, err := uuid.Parse(stringUUID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse UUID: %v", stringUUID)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	me, err := resources.PostgresDb.GetUserById(request.Context(), userID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not find user: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	chatsWithMeAsCreator, err := resources.PostgresDb.GetChatsByCreatorId(request.Context(), userID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not find chats: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	chatsWithMeAsFriend, err := resources.PostgresDb.GetChatsByFriendId(request.Context(), userID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not find chats: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	chats := []Chat{}

	for _, chat := range chatsWithMeAsCreator {
		chats = append(chats, chatsWithMeAsCreatorToChat(me, chat))
	}

	for _, chat := range chatsWithMeAsFriend {
		chats = append(chats, chatsWithMeAsFriendToChat(me, chat))
	}

	helpers.ResponseJson(writer, 200, Chats{Data: chats})
}
