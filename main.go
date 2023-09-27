package main

import (
	"database/sql"
	"github.com/dnieln7/just-chatting/internal/env"
	"log"
	"net/http"
	"time"

	"github.com/dnieln7/just-chatting/internal/database/db"
	"github.com/dnieln7/just-chatting/internal/handler/chat"
	"github.com/dnieln7/just-chatting/internal/handler/message"
	"github.com/dnieln7/just-chatting/internal/handler/user"
	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/dnieln7/just-chatting/internal/server/chatserver"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	properties := env.GetEnvProperties()

	resources := &server.Resources{
		PostgresDb: buildDatabase(properties),
	}

	router := buildRouter(resources)

	chatServer := buildChatServer(resources, router)
	httpServer := buildHttpServer(properties, router)

	chatServer.ListenAndServe()
	err := httpServer.ListenAndServe()

	if err != nil {
		log.Fatal("Could not start server", err)
	}
}

func buildDatabase(properties *env.EvnProperties) *db.Queries {
	connection, err := sql.Open("postgres", properties.PostgresUrl)

	if err != nil {
		log.Fatal("Could not connect to database")
	}

	queries := db.New(connection)

	return queries
}

func buildChatServer(resources *server.Resources, router *mux.Router) *chatserver.ChatServer {
	chatServer := &chatserver.ChatServer{
		Resources:         resources,
		IncomingMessages:  make(chan chatserver.IncomingMessage),
		ConnectionUpdates: make(chan chatserver.ConnectionUpdate),
	}

	router.HandleFunc("/users/{user_id}/connect/{chat_id}", chatServer.UpgraderHandler).
		Methods("GET")

	return chatServer
}

func buildRouter(resources *server.Resources) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/signup", resources.HttpHandler(user.PostUserHandler)).
		Methods("POST")

	router.HandleFunc("/login", resources.HttpHandler(user.GetUserByEmailHandler)).
		Methods("POST")

	router.HandleFunc("/users/{id}/chats", resources.HttpHandler(chat.GetChatsByParticipantIdHandler)).
		Methods("GET")

	router.HandleFunc("/messages", resources.HttpHandler(message.PostMessageHandler)).
		Methods("POST")

	router.HandleFunc("/chats", resources.HttpHandler(chat.PostChatHandler)).
		Methods("POST")

	router.HandleFunc("/chats/{id}/messages", resources.HttpHandler(message.GetMessagesByChatIdHandler)).
		Queries("page", "{page:[0-9]+}").Methods("GET")

	return router
}
func buildHttpServer(properties *env.EvnProperties, router *mux.Router) *http.Server {
	httpServer := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:" + properties.Port,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	return httpServer
}
