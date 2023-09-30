package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dnieln7/just-chatting/internal/database/db"
	"github.com/dnieln7/just-chatting/internal/helpers"
	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/google/uuid"
)

type PostChatBody struct {
	Participants []uuid.UUID `json:"participants"`
}

func PostChatHandler(writer http.ResponseWriter, request *http.Request, resources *server.Resources) {
	decoder := json.NewDecoder(request.Body)
	body := PostChatBody{}
	err := decoder.Decode(&body)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse JSON: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	existingDbChat, err := resources.PostgresDb.GetChatWithParticipants(request.Context(), body.Participants)

	if err != nil {
		errMessage := fmt.Sprintf("Could not create chat: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	if len(existingDbChat.Participants) != 0 {
		helpers.ResponseJson(writer, 409, dbChatToChat(existingDbChat))
		return
	}

	dbChat, err := resources.PostgresDb.CreateChat(request.Context(), db.CreateChatParams{
		ID:           uuid.New(),
		Participants: body.Participants,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})

	if err != nil {
		errMessage := fmt.Sprintf("Could not create chat: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
	} else {
		helpers.ResponseJson(writer, 201, dbChatToChat(dbChat))
	}
}
