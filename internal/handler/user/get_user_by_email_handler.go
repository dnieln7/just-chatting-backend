package user

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/dnieln7/just-chatting/internal/helpers"
	"github.com/dnieln7/just-chatting/internal/server"
)

func GetUserByEmailHandler(writer http.ResponseWriter, request *http.Request, resources *server.Resources) {
	vars := mux.Vars(request)
	email := vars["email"]

	dbUser, err := resources.PostgresDb.GetUserByEmail(request.Context(), email)

	if err != nil {
		errMessage := fmt.Sprintf("Could not get user: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	helpers.ResponseJson(writer, 200, dbUserToUser(dbUser))
}
