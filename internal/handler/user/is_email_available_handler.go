package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dnieln7/just-chatting/internal/helpers"
	"github.com/dnieln7/just-chatting/internal/server"
)

type IsEmailAvailableBody struct {
	Email string `json:"email"`
}

func IsEmailAvailableHandler(writer http.ResponseWriter, request *http.Request, resources *server.Resources) {
	decoder := json.NewDecoder(request.Body)
	body := IsEmailAvailableBody{}
	err := decoder.Decode(&body)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse JSON: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	count, err := resources.PostgresDb.CountUsersWithEmail(request.Context(), body.Email)

	if err != nil {
		errMessage := fmt.Sprintf("Could not check email availability: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	if count > 0 {
		helpers.ResponseJsonError(writer, 409, "Email is not available")
	} else {
		helpers.ResponseNoContent(writer)
	}
}
