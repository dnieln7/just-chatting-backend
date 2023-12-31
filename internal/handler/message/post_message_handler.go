package message

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

type PostMessageBody struct {
	ChatID  uuid.UUID `json:"chat_id"`
	UserID  uuid.UUID `json:"user_id"`
	Message string    `json:"message"`
}

// TODO: only participants should be able to send messages
func PostMessageHandler(writer http.ResponseWriter, request *http.Request, resources *server.Resources) {
	decoder := json.NewDecoder(request.Body)
	body := PostMessageBody{}
	err := decoder.Decode(&body)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse JSON: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	dbMessage, err := resources.PostgresDb.CreateMessage(request.Context(), db.CreateMessageParams{
		ID:        uuid.New(),
		ChatID:    body.ChatID,
		UserID:    body.UserID,
		Message:   body.Message,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		errMessage := fmt.Sprintf("Could not create message: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
	} else {
		helpers.ResponseJson(writer, 201, DBMessageToMessage(dbMessage))
	}
}
