package server

import (
	"database/sql"
	"net/http"

	"github.com/dnieln7/just-chatting/internal/database/db"
)

type Resources struct {
	ConnectionDb *sql.DB
	PostgresDb   *db.Queries
}

type ResourcesHandlerFunc func(writer http.ResponseWriter, request *http.Request, resources *Resources)

func (resources *Resources) HttpHandler(handlerFunc ResourcesHandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		handlerFunc(writer, request, resources)
	}
}
