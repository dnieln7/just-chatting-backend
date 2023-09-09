package model

import "github.com/gorilla/websocket"

type ConnectionUpdate struct {
	Connection *websocket.Conn
	Register bool
}
