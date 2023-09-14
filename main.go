package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dnieln7/just-chatting/internal/database/db"
	"github.com/dnieln7/just-chatting/internal/handler/chat"
	"github.com/dnieln7/just-chatting/internal/handler/message"
	"github.com/dnieln7/just-chatting/internal/handler/user"
	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/dnieln7/just-chatting/internal/server/chatserver"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var upgrader = websocket.Upgrader{}

func main() {
	godotenv.Load()

	serverResources := server.ServerResources{
		PostgresDb: setupDatabase(),
	}

	chatServer := chatserver.ChatServer{
		Resources:         &serverResources,
		IncomingMessages:  make(chan chatserver.IncomingMessage),
		ConnectionUpdates: make(chan chatserver.ConnectionUpdate),
	}

	router := mux.NewRouter()

	router.HandleFunc("/signup", serverResources.WithResources(user.PostUserHandler)).
		Methods("POST")

	router.HandleFunc("/login", serverResources.WithResources(user.GetUserByEmailHandler)).
		Methods("POST")

	router.HandleFunc("/users/{id}/chats", serverResources.WithResources(chat.GetChatsByParticipantIdHandler)).
		Methods("GET")

	router.HandleFunc("/messages", serverResources.WithResources(message.PostMessageHandler)).
		Methods("POST")

	router.HandleFunc("/chats", serverResources.WithResources(chat.PostChatHandler)).
		Methods("POST")

	router.HandleFunc("/chats/{id}/messages", serverResources.WithResources(message.GetMessagesByChatIdHandler)).
		Queries("page", "{page:[0-9]+}").Methods("GET")

	router.HandleFunc("/users/{user_id}/connect/{chat_id}", chatServer.UpgraderHandler).
		Methods("GET")

	chatServer.ListenAndServe()
	setupServer(router)
}

func setupDatabase() *db.Queries {
	dbUrl := os.Getenv("DB_URL")

	if dbUrl == "" {
		log.Fatal("DB_URL not found")
	}

	connection, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Fatal("Could not connect to database")
	}

	queries := db.New(connection)

	return queries
}

func setupServer(router *mux.Router) {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("PORT not found")
	} else {
		log.Println("Starting server on port: ", port, "...")
	}

	server := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:" + port,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal("Could not start server", err)
	}
}
