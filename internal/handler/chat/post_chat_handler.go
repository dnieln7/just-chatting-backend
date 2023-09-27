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
	Participants []string `json:"participants"`
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

	participants := []uuid.UUID{}

	for _, stringUUID := range body.Participants {
		participant, err := uuid.Parse(stringUUID)

		if err != nil {
			errMessage := fmt.Sprintf("Could not parse UUID: %v", stringUUID)
			helpers.ResponseJsonError(writer, 400, errMessage)
			return
		}

		participants = append(participants, participant)
	}

	existingDbChat, err := resources.PostgresDb.GetChatWithParticipants(request.Context(), participants)

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
		Participants: participants,
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
