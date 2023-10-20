package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dnieln7/just-chatting/internal/helpers"
	"github.com/dnieln7/just-chatting/internal/server"
)

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(writer http.ResponseWriter, request *http.Request, resources *server.Resources) {
	decoder := json.NewDecoder(request.Body)
	body := LoginBody{}
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
