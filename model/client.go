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
	Saved	SavedChat
}

func (client *Client) Read() {
	defer client.Socket.Close()
	for {
		var chat *model.Chat
		err := client.Socket.ReadJSON(&chat)
		if err != nil {
			log.Print(err)
			return
		}

		chat.Name = client.User.Name
		chat.Channel = client.Channel
		chat.User = client.User.ID
		chat.Timestamp = time.Now()
		
		client.Room.Forward <- chat

		savedChat := &SavedChat{
			Chat: chat
		}

		*client.Saved <- savedChat
	}
}

func (client *Client) Write() {
	defer client.Socket.Close()
	for chat := range client.Send {
		err := client.Socket.WriteJSON(chat)
		if err != nil {
			return
		}
	}
}