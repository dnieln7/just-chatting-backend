package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dnieln7/just-chatting/internal/database/db"
	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/google/uuid"
)

type PostUserBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func PostUserHandler(writer http.ResponseWriter, request *http.Request, resources *server.ServerResources) {
	decoder := json.NewDecoder(request.Body)
	body := PostUserBody{}
	err := decoder.Decode(&body)

	log.Println(body)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse JSON: %v", err)
		server.ResponseJsonError(writer, 400, errMessage)
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
		server.ResponseJsonError(writer, 400, errMessage)
	} else {
		server.ResponseJson(writer, 201, dbUserToUser(dbUser))
	}
}
