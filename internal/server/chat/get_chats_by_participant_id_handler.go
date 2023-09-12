package chat

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/google/uuid"
)

func GetChatsByParticipantIdHandler(writer http.ResponseWriter, request *http.Request, resources *server.ServerResources) {
	path := strings.Split(request.URL.Path, "/")
	stringUUID := path[1]

	participant, err := uuid.Parse(stringUUID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse UUID: %v", stringUUID)
		server.ResponseJsonError(writer, 400, errMessage)
		return
	}

	dbChats, err := resources.PostgresDb.GetChatsByParticipantId(request.Context(), []uuid.UUID{participant})

	if err != nil {
		errMessage := fmt.Sprintf("Could not create chat: %v", err)
		server.ResponseJsonError(writer, 400, errMessage)
	} else {
		chats := []Chat{}

		for _, dbChat := range dbChats {
			chats = append(chats, dbChatToChat(dbChat))
		}

		server.ResponseJson(writer, 201, Chats{Data: chats})
	}
}
