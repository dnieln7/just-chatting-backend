package main

import (
	"log"
	"net/http"

	"github.com/dnieln7/just-chatting/internal/model"
	"github.com/dnieln7/just-chatting/internal/server/chat"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {
	chatServer := setUpChatServer()

	http.HandleFunc("/chat", chatHandler(chatServer))

	log.Println("Starting server...")
	err := http.ListenAndServe(":4444", nil)

	if err != nil {
		log.Fatal("Could not start server", err)
	}
}

func setUpChatServer() *chat.ChatServer {
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
