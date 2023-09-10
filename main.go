package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/dnieln7/just-chatting/internal/database/db"
	"github.com/dnieln7/just-chatting/internal/model"
	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/dnieln7/just-chatting/internal/server/chat"
	"github.com/dnieln7/just-chatting/internal/server/user"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var upgrader = websocket.Upgrader{}

func main() {
	godotenv.Load()

	serverResources := server.ServerResources {
		PostgresDb: setupDatabase(),	
	}

	chatServer := setupChatServer()

	http.HandleFunc("/chat", chatHandler(chatServer))
	http.HandleFunc("/signup", serverResources.WithResources(user.PostUserHandler))
	http.HandleFunc("/login", serverResources.WithResources(user.GetUserByEmailHandler))

	setupServer()
}

func setupChatServer() *chat.ChatServer {
	chatServer := chat.ChatServer{
		Messages:          make(chan []byte),
		ConnectionUpdates: make(chan model.ConnectionUpdate),
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

func setupServer()  {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("PORT not found")
	}

	log.Println("Starting server on port: ", port, "...")

	err := http.ListenAndServe(":" + port, nil)

	if err != nil {
		log.Fatal("Could not start server", err)
	}
}

func chatHandler(chat *chat.ChatServer) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		connection, err := upgrader.Upgrade(writer, request, nil)

		if err != nil {
			log.Println("Could not upgrade request: ", err)
			return
		}

		chat.ConnectionUpdates <- model.ConnectionUpdate{
			Connection: connection,
			Register:   true,
		}
	}
}
