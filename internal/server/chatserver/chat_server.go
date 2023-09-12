package chatserver

import (
	"fmt"
	"log"
	"github.com/gorilla/websocket"
)

type ChatServer struct {
	connections       []*websocket.Conn
	Messages          chan []byte
	ConnectionUpdates chan ConnectionUpdate
}

func (chat *ChatServer) AddConnection(connection *websocket.Conn) {
	go func() {
		for {
			messageType, message, err := connection.ReadMessage()

			if err != nil {
				closeErr, ok := err.(*websocket.CloseError)

				if ok {
					log.Println("Close frame received, clossing...", closeErr)

					chat.ConnectionUpdates <- ConnectionUpdate{
						Connection: connection,
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

			chat.Messages <- message
		}
	}()
}

func (chat *ChatServer) RemoveConnection(connection *websocket.Conn) {

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

		log.Println("Removing connection... ", connection.LocalAddr(), " at index: ", index)

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

func (chat *ChatServer) WriteMessage(message []byte) {
	for _, connection := range chat.connections {
		err := connection.WriteMessage(websocket.TextMessage, message)

		if err != nil {
			log.Println("Error writing message:", err)
		} else {
			messageText := fmt.Sprintf("%s", message)

			log.Println("Message sent: ", messageText)
		}
	}
}
