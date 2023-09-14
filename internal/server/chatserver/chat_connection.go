package chatserver

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ChatConnection struct {
	Connection *websocket.Conn
	ChatID uuid.UUID
}
