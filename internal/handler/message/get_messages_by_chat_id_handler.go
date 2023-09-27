package message

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dnieln7/just-chatting/internal/database/db"
	"github.com/dnieln7/just-chatting/internal/helpers"
	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetMessagesByChatIdHandler(writer http.ResponseWriter, request *http.Request, resources *server.Resources) {
	vars := mux.Vars(request)
	stringUUID := vars["id"]
	stringPage := vars["page"]

	chatID, err := uuid.Parse(stringUUID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse UUID: %v", stringUUID)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	page, _ := strconv.Atoi(stringPage)

	if page < 1 {
		helpers.ResponseJsonError(writer, 400, "Page is invalid")
		return
	}

	dbMessages, err := resources.PostgresDb.GetMessagesByChatIdLazy(request.Context(), db.GetMessagesByChatIdLazyParams{
		ChatID: chatID,
		Limit:  PAGE_SIZE,
		Offset: int32((page - 1) * PAGE_SIZE),
	})

	var hasNextPage bool

	if len(dbMessages) < PAGE_SIZE {
		hasNextPage = false
	} else {
		hasNextPage = true
	}

	if err != nil {
		errMessage := fmt.Sprintf("Could not get messages: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
	} else {
		messages := []Message{}

		for _, dbMessage := range dbMessages {
			messages = append(messages, dbMessageToMessage(dbMessage))
		}

		helpers.ResponseJson(writer, 200, Messages{
			Data:        messages,
			CurrentPage: page,
			HasNextPage: hasNextPage,
		})
	}
}

const PAGE_SIZE = 100
