package model

import (
	"github.com/gorilla/websocket"
	"github.com/ivohutasoit/alira/model/domain"
)

// Client represents a user connect to a room, one user may have many devices to chat,
// so it should not be the same as user
type Client struct {
	Channel string
	Socket  *websocket.Conn
	Send    chan *domain.Chat
	Room    Room
	User    domain.User
}
