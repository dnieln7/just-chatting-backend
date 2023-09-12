package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dnieln7/just-chatting/internal/helpers"
	"github.com/dnieln7/just-chatting/internal/server"
)

type GetUserByEmailBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetUserByEmailHandler(writer http.ResponseWriter, request *http.Request, resources *server.ServerResources) {
	decoder := json.NewDecoder(request.Body)
	body := GetUserByEmailBody{}
	err := decoder.Decode(&body)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse JSON: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	dbUser, err := resources.PostgresDb.GetUserByEmail(request.Context(), body.Email)

	if err != nil {
		errMessage := fmt.Sprintf("Could not get user: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	if dbUser.Password == body.Password {
		helpers.ResponseJson(writer, 200, dbUserToUser(dbUser))
	} else {
		helpers.ResponseJsonError(writer, 401, "Wrong password")
	}
}
