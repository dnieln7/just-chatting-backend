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
	CreatorID uuid.UUID `json:"creator_id"`
	FriendID  uuid.UUID `json:"friend_id"`
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

	existingDbChat, err := resources.PostgresDb.GetChatWithCreatorAndFriend(request.Context(), db.GetChatWithCreatorAndFriendParams{
		CreatorID: body.CreatorID,
		FriendID:  body.FriendID,
	})

	if err != nil {
		errMessage := fmt.Sprintf("Could not create chat: %v", err)

		if errMessage != "Could not create chat: sql: no rows in result set" {
			helpers.ResponseJsonError(writer, 400, errMessage)
			return
		}
	}

	existsWithMeAsCreator := existingDbChat.CreatorID == body.CreatorID && existingDbChat.FriendID == body.FriendID
	existsWithMeAsFriend := existingDbChat.CreatorID == body.FriendID && existingDbChat.FriendID == body.CreatorID

	if existsWithMeAsCreator {
		errMessage, creator, friend := getCreatorAndFriend(resources.PostgresDb, request.Context(), body.CreatorID, body.FriendID)

		if errMessage != "" {
			helpers.ResponseJsonError(writer, 400, errMessage)
			return
		}

		chat := Chat{
			ID: existingDbChat.ID.String(),
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
			CreatedAt: existingDbChat.CreatedAt,
			UpdatedAt: existingDbChat.UpdatedAt,
		}

		helpers.ResponseJson(writer, 409, chat)
		return
	}

	if existsWithMeAsFriend {
		errMessage, creator, friend := getCreatorAndFriend(resources.PostgresDb, request.Context(), body.FriendID, body.CreatorID)

		if errMessage != "" {
			helpers.ResponseJsonError(writer, 400, errMessage)
			return
		}

		chat := Chat{
			ID: existingDbChat.ID.String(),
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
			CreatedAt: existingDbChat.CreatedAt,
			UpdatedAt: existingDbChat.UpdatedAt,
		}

		helpers.ResponseJson(writer, 409, chat)
		return
	}

	dbChat, err := resources.PostgresDb.CreateChat(request.Context(), db.CreateChatParams{
		ID:        uuid.New(),
		CreatorID: body.CreatorID,
		FriendID:  body.FriendID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		errMessage := fmt.Sprintf("Could not create chat: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
	} else {
		errMessage, creator, friend := getCreatorAndFriend(resources.PostgresDb, request.Context(), body.CreatorID, body.FriendID)

		if errMessage != "" {
			helpers.ResponseJsonError(writer, 400, errMessage)
			return
		}

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

		helpers.ResponseJson(writer, 201, chat)
		helpers.ResponseOK(writer)
	}
}
