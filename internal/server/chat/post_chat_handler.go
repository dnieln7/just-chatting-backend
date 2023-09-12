package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dnieln7/just-chatting/internal/database/db"
	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/google/uuid"
)

type PostChatBody struct {
	Participants []string `json:"participants"`
}

func PostChatHandler(writer http.ResponseWriter, request *http.Request, resources *server.ServerResources) {
	decoder := json.NewDecoder(request.Body)
	body := PostChatBody{}
	err := decoder.Decode(&body)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse JSON: %v", err)
		server.ResponseJsonError(writer, 400, errMessage)
		return
	}

	participants := []uuid.UUID{}

	for _, stringUUID := range body.Participants {
		participant, err := uuid.Parse(stringUUID)

		if err != nil {
			errMessage := fmt.Sprintf("Could not parse UUID: %v", stringUUID)
			server.ResponseJsonError(writer, 400, errMessage)
			return
		}

		participants = append(participants, participant)
	}

	dbChat, err := resources.PostgresDb.CreateChat(request.Context(), db.CreateChatParams{
		ID:           uuid.New(),
		Participants: participants,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})

	if err != nil {
		errMessage := fmt.Sprintf("Could not create chat: %v", err)
		server.ResponseJsonError(writer, 400, errMessage)
	} else {
		server.ResponseJson(writer, 201, dbChatToChat(dbChat))
	}
}
