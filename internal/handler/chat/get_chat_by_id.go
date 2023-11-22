package chat

import (
	"fmt"
	"github.com/dnieln7/just-chatting/internal/helpers"
	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

func GetChatByIdHandler(writer http.ResponseWriter, request *http.Request, resources *server.Resources) {
	vars := mux.Vars(request)
	stringUserUUID := vars["user_id"]
	stringChatUUID := vars["chat_id"]

	userID, err := uuid.Parse(stringUserUUID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse UUID: %v", stringUserUUID)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	chatID, err := uuid.Parse(stringChatUUID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse UUID: %v", stringChatUUID)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	dbChat, err := resources.PostgresDb.GetChatById(request.Context(), chatID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not find chat: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	iAmCreator := dbChat.CreatorID == userID
	iAmFriend := dbChat.FriendID == userID

	if !iAmCreator && !iAmFriend {
		helpers.ResponseJsonError(writer, 401, "You are not a participant of this chat")
		return
	}

	errMessage, creator, friend := getCreatorAndFriend(resources.PostgresDb, request.Context(), dbChat.CreatorID, dbChat.FriendID)

	if errMessage != "" {
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	if iAmCreator {
		chat := Chat{
			ID: dbChat.ID.String(),
			Me: Participant{
				ID:       creator.ID.String(),
				Email:    creator.Email,
				Username: creator.Username,
			},
			Creator: Participant{
				ID:       creator.ID.String(),
				Email:    creator.Email,
				Username: creator.Username,
			},
			Friend: Participant{
				ID:       friend.ID.String(),
				Email:    friend.Email,
				Username: friend.Username,
			},
			CreatedAt: dbChat.CreatedAt,
			UpdatedAt: dbChat.UpdatedAt,
		}

		helpers.ResponseJson(writer, 200, chat)
	} else {
		chat := Chat{
			ID: dbChat.ID.String(),
			Me: Participant{
				ID:       friend.ID.String(),
				Email:    friend.Email,
				Username: friend.Username,
			},
			Creator: Participant{
				ID:       creator.ID.String(),
				Email:    creator.Email,
				Username: creator.Username,
			},
			Friend: Participant{
				ID:       creator.ID.String(),
				Email:    creator.Email,
				Username: creator.Username,
			},
			CreatedAt: dbChat.CreatedAt,
			UpdatedAt: dbChat.UpdatedAt,
		}

		helpers.ResponseJson(writer, 200, chat)
	}
}
