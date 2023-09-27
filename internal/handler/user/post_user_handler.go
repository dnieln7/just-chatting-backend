package user

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

type PostUserBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func PostUserHandler(writer http.ResponseWriter, request *http.Request, resources *server.Resources) {
	decoder := json.NewDecoder(request.Body)
	body := PostUserBody{}
	err := decoder.Decode(&body)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse JSON: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	dbUser, err := resources.PostgresDb.CreateUser(request.Context(), db.CreateUserParams{
		ID:        uuid.New(),
		Email:     body.Email,
		Password:  body.Password,
		Username:  body.Username,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		errMessage := fmt.Sprintf("Could not create user: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
	} else {
		helpers.ResponseJson(writer, 201, dbUserToUser(dbUser))
	}
}
