package chatserver

import "github.com/google/uuid"

type IncomingMessage struct {
	Message []byte
	UserID  uuid.UUID
	ChatID  uuid.UUID
}
