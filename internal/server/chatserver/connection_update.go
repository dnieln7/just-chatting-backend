package chatserver

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ConnectionUpdate struct {
	Connection *websocket.Conn
	UserID uuid.UUID
	ChatID uuid.UUID
	Register bool
}
