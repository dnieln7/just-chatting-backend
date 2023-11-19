package chatserver

import (
	"context"
	"fmt"
	"github.com/dnieln7/just-chatting/internal/database/db"
	"log"
	"net/http"
	"time"

	"github.com/dnieln7/just-chatting/internal/helpers"
	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type ChatServer struct {
	connections       []ChatConnection
	Resources         *server.Resources
	IncomingMessages  chan IncomingMessage
	ConnectionUpdates chan ConnectionUpdate
}

func (chat *ChatServer) ListenAndServe() {
	go func() {
		for {
			select {
			case incomingMessage := <-chat.IncomingMessages:
				chat.BroadcastMessage(incomingMessage)
				SaveMessage(chat.Resources.PostgresDb, incomingMessage)
			case connectionUpdate := <-chat.ConnectionUpdates:
				if connectionUpdate.Register {
					chat.RegisterConnection(connectionUpdate)
				} else {
					chat.UnregisterConnection(connectionUpdate.Connection)
				}
			}
		}
	}()
}

func (chat *ChatServer) UpgraderHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	stringUserUUID := vars["user_id"]
	stringChatUUID := vars["chat_id"]

	chatID, err := uuid.Parse(stringChatUUID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse UUID: %v", stringChatUUID)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	userID, err := uuid.Parse(stringUserUUID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not parse UUID: %v", stringUserUUID)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	}

	dbChat, err := chat.Resources.PostgresDb.GetChatById(request.Context(), chatID)

	if err != nil {
		errMessage := fmt.Sprintf("Could not find chats: %v", err)
		helpers.ResponseJsonError(writer, 400, errMessage)
		return
	} else {
		log.Printf("Chat %v found checking participants...\n", chatID)
	}

	isUserParticipant := helpers.ContainsUUID([]uuid.UUID{dbChat.CreatorID, dbChat.FriendID}, userID)

	if !isUserParticipant {
		helpers.ResponseJsonError(writer, 401, "You are not a participant in this chat")
		return
	} else {
		log.Printf("Chat %v has participant %v upgrading...\n", chatID, userID)
	}

	connection, err := upgrader.Upgrade(writer, request, nil)

	if err != nil {
		log.Println("Could not upgrade request: ", err)
		return
	}

	chat.ConnectionUpdates <- ConnectionUpdate{
		Connection: connection,
		UserID:     userID,
		ChatID:     chatID,
		Register:   true,
	}
}

func (chat *ChatServer) RegisterConnection(connectionUpdate ConnectionUpdate) {
	log.Printf("Registering connection %v ...\n", connectionUpdate.Connection.RemoteAddr())

	chat.connections = append(chat.connections, ChatConnection{
		Connection: connectionUpdate.Connection,
		UserID:     connectionUpdate.UserID,
		ChatID:     connectionUpdate.ChatID,
	})

	go func() {
		for {
			messageType, message, err := connectionUpdate.Connection.ReadMessage()

			if err != nil {
				closeErr, ok := err.(*websocket.CloseError)

				if ok {
					log.Println("Close frame received, closing...", closeErr)

					chat.ConnectionUpdates <- ConnectionUpdate{
						Connection: connectionUpdate.Connection,
						UserID:     connectionUpdate.UserID,
						ChatID:     connectionUpdate.ChatID,
						Register:   false,
					}

					break
				} else {
					log.Println("Error reading message: ", err)
					continue
				}
			}

			messageText := fmt.Sprintf("%s", message)

			log.Println("Message received: ", messageText, " with type: ", messageType)

			chat.IncomingMessages <- IncomingMessage{
				Message: message,
				UserID:  connectionUpdate.UserID,
				ChatID:  connectionUpdate.ChatID,
			}
		}
	}()
}

func (chat *ChatServer) UnregisterConnection(connection *websocket.Conn) {
	log.Println("Unregistering connection...")

	var index = -1
	var last = len(chat.connections) - 1

	if last == 0 {
		log.Printf("Removing connection %v at index 0...\n", connection.RemoteAddr())

		chat.connections = []ChatConnection{}
	} else {
		for i, conn := range chat.connections {
			if conn.Connection == connection {
				index = i
				break
			}
		}

		log.Printf("Removing connection %v at index %d...\n", connection.RemoteAddr(), index)

		if index != -1 {
			if index != last {
				chat.connections[index] = chat.connections[last]
			}

			chat.connections = chat.connections[:last]
		}
	}

	log.Printf("Connection removed, remaining connections %d\n", len(chat.connections))

	err := connection.Close()

	if err != nil {
		log.Printf("Error closing connection %v: %v\n", connection.RemoteAddr(), err)
	}
}

func (chat *ChatServer) BroadcastMessage(incomingMessage IncomingMessage) {
	dbChat, err := chat.Resources.PostgresDb.GetChatById(context.Background(), incomingMessage.ChatID)
	participants := []uuid.UUID{dbChat.CreatorID, dbChat.FriendID}

	if err != nil {
		log.Printf("Error finding participants fo chat %v: %v\n", incomingMessage.ChatID, err)
	}

	participants = helpers.RemoveUUID(participants, incomingMessage.UserID)

	log.Printf("Writing message to chat %v...\n", incomingMessage.ChatID)

	for _, conn := range chat.connections {
		if conn.ChatID == incomingMessage.ChatID && conn.UserID != incomingMessage.UserID {
			err := conn.Connection.WriteMessage(websocket.TextMessage, incomingMessage.Message)

			if err != nil {
				log.Printf("Error writing message to chat %v\n", incomingMessage.ChatID)
			} else {
				log.Printf("A message was sent to connection %v of chat %v\n", conn.Connection.RemoteAddr(), incomingMessage.ChatID)

				participants = helpers.RemoveUUID(participants, conn.UserID)
			}
		}
	}

	log.Printf("Participants of chat %v without a connection or with a write error %d\n", incomingMessage.ChatID, len(participants))
}

func SaveMessage(postgresDb *db.Queries, incomingMessage IncomingMessage) {
	messageText := fmt.Sprintf("%s", incomingMessage.Message)

	_, err := postgresDb.CreateMessage(context.Background(), db.CreateMessageParams{
		ID:        uuid.New(),
		ChatID:    incomingMessage.ChatID,
		UserID:    incomingMessage.UserID,
		Message:   messageText,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Printf("Error saving message to chat %v: %v\n", incomingMessage.ChatID, err)
	} else {
		log.Printf("Saved message to chat %v\n", incomingMessage.ChatID)
	}
}
