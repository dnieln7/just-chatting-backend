package chat

import (
	"fmt"
	"net/http"

	"github.com/dnieln7/just-chatting/internal/helpers"
	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetChatsByParticipantIdHandler(writer http.ResponseWriter, request *http.Request, resources *server.Resources) {
	vars := mux.Vars(request)
	stringUUID := vars["id"]

	participant, err := uuid.Parse(stringUUID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse UUID: %v", stringUUID)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	dbChats, err := resources.PostgresDb.GetChatsByParticipantId(request.Context(), []uuid.UUID{participant})

	if err != nil {
		errMessage := fmt.Sprintf("Could not find chats: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
	} else {
		chats := []Chat{}

		for _, dbChat := range dbChats {
			chats = append(chats, dbChatToChat(dbChat))
		}

		helpers.ResponseJson(writer, 200, Chats{Data: chats})
	}
}
