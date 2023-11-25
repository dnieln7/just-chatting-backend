package message

import (
	"fmt"
	"github.com/dnieln7/just-chatting/internal/helpers"
	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

func GetMessagesByChatIdEagerHandler(writer http.ResponseWriter, request *http.Request, resources *server.Resources) {
	vars := mux.Vars(request)
	stringUUID := vars["id"]

	chatID, err := uuid.Parse(stringUUID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse UUID: %v", stringUUID)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	dbMessages, err := resources.PostgresDb.GetMessagesByChatIdEager(request.Context(), chatID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not get messages: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
	} else {
		messages := []Message{}

		for _, dbMessage := range dbMessages {
			messages = append(messages, DBMessageToMessage(dbMessage))
		}

		helpers.ResponseJson(writer, 200, Messages{
			Data:        messages,
			CurrentPage: 1,
			HasNextPage: false,
		})
	}
}
