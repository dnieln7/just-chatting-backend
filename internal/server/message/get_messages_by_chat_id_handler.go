package message

import (
	"fmt"
	"net/http"

	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetMessagesByChatIdHandler(writer http.ResponseWriter, request *http.Request, resources *server.ServerResources) {
	vars := mux.Vars(request)
	stringUUID := vars["id"]

	chatID, err := uuid.Parse(stringUUID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse UUID: %v", stringUUID)
		server.ResponseJsonError(writer, 400, errMessage)
		return
	}

	dbMessages, err := resources.PostgresDb.GetMessagesByChatId(request.Context(), chatID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not create chat: %v", err)
		server.ResponseJsonError(writer, 400, errMessage)
	} else {
		messages := []Message{}

		for _, dbMessage := range dbMessages {
			messages = append(messages, dbMessageToMessage(dbMessage))
		}

		server.ResponseJson(writer, 200, Messages{Data: messages})
	}
}
