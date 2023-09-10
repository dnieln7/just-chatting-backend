package server

import (
	"net/http"

	"github.com/dnieln7/just-chatting/internal/database/db"
)

type ServerResources struct {
	PostgresDb *db.Queries
}

type ResourcesHandlerFunc func(writer http.ResponseWriter, request *http.Request, resources *ServerResources)

func (serverResources *ServerResources) WithResources(handlerFunc ResourcesHandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		handlerFunc(writer, request, serverResources)
	}
}
