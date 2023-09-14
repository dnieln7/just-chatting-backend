package chatserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dnieln7/just-chatting/internal/helpers"
	"github.com/dnieln7/just-chatting/internal/server"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type ChatServer struct {
	connections       []*websocket.Conn
	Resources         *server.ServerResources
	IncomingMessages  chan IncomingMessage
	ConnectionUpdates chan ConnectionUpdate
}

func (chat *ChatServer) ListenAndServe() {
	go func() {
		for {
			select {
			case incomingMessage := <-chat.IncomingMessages:
				chat.WriteMessage(incomingMessage)
			case connectionUpdate := <-chat.ConnectionUpdates:
				if connectionUpdate.Register {
					chat.AddConnectionUpdate(connectionUpdate)
				} else {
					chat.RemoveConnection(connectionUpdate.Connection)
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

	isUserParticipant := helpers.ContainsUUID(dbChat.Participants, userID)

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

func (chat *ChatServer) AddConnectionUpdate(connectionUpdate ConnectionUpdate) {
	log.Printf("Registering connection %v ...\n", connectionUpdate.Connection.RemoteAddr())

	chat.connections = append(chat.connections, connectionUpdate.Connection)

	go func() {
		for {
			messageType, message, err := connectionUpdate.Connection.ReadMessage()

			if err != nil {
				closeErr, ok := err.(*websocket.CloseError)

				if ok {
					log.Println("Close frame received, clossing...", closeErr)

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

func (chat *ChatServer) RemoveConnection(connection *websocket.Conn) {
	log.Printf("Unregistering connection %v ...\n", connection.RemoteAddr())

	var index = -1
	var last = len(chat.connections) - 1

	if last == 0 {
		log.Println("Cleaning connections... ")

		chat.connections = []*websocket.Conn{}
	} else {
		for i, conn := range chat.connections {
			if conn == connection {
				index = i
				break
			}
		}

		log.Println("Removing connection... ", connection.RemoteAddr(), " at index: ", index)

		if index != -1 {
			if index != last {
				chat.connections[index] = chat.connections[last]
			}

			chat.connections = chat.connections[:last]
		}

		log.Println("Connection removed, remaining connections: ", len(chat.connections))
	}

	connection.Close()
}

func (chat *ChatServer) WriteMessage(incomingMessage IncomingMessage) {
	log.Println("WriteMessage...")
	for _, connection := range chat.connections {
		err := connection.WriteMessage(websocket.TextMessage, incomingMessage.Message)

		if err != nil {
			log.Println("Error writing message:", err)
		} else {
			messageText := fmt.Sprintf("%s", incomingMessage.Message)

			log.Println("Message sent: ", messageText)
		}
	}
}
