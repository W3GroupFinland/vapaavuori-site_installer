package web_models

import (
	"github.com/gorilla/websocket"
)

// TODO: Implement date..
func NewStatusMessage(msg string) *StatusMessage {
	return &StatusMessage{Message: msg}
}

// TODO: Implement date..
func NewConnStatusMessage(conn *websocket.Conn, msg string) *ConnStatusMessage {
	csm := ConnStatusMessage{Connection: conn}
	csm.Message = msg
	return &csm
}

type StatusMessage struct {
	Date    string
	Message string
}

type ConnStatusMessage struct {
	StatusMessage
	Connection *websocket.Conn
}
