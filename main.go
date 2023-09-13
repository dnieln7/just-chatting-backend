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

	// http.HandleFunc("/chat", chatHandler(chatServer))

	// chatServer := setupChatServer()
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

func setupChatServer() *chatserver.ChatServer {
	chatServer := chatserver.ChatServer{
		Messages:          make(chan []byte),
		ConnectionUpdates: make(chan chatserver.ConnectionUpdate),
	}

	go func() {
		for {
			select {
			case message := <-chatServer.Messages:
				chatServer.WriteMessage(message)
			case connectionUpdate := <-chatServer.ConnectionUpdates:
				if connectionUpdate.Register {
					chatServer.AddConnection(connectionUpdate.Connection)
				} else {
					chatServer.RemoveConnection(connectionUpdate.Connection)
				}
			}
		}
	}()

	return &chatServer
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

func chatHandler(chat *chatserver.ChatServer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		connection, err := upgrader.Upgrade(writer, request, nil)

		if err != nil {
			log.Println("Could not upgrade request: ", err)
			return
		}

		chat.ConnectionUpdates <- chatserver.ConnectionUpdate{
			Connection: connection,
			Register:   true,
		}
	}
}
