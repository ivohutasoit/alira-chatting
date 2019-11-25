package model

import (
	"github.com/ivohutasoit/alira/model/domain"
)

// Room represents a room to chat
type Room struct {
	Forward chan *domain.Chat
	Join    chan *Client
	Leave   chan *Client
	Clients map[*Client]bool
}
