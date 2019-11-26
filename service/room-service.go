package service

import (
	"github.com/gorilla/websocket"
	"github.com/ivohutasoit/alira-chatting/model"
	"github.com/ivohutasoit/alira/model/domain"
	"log"
	"net/http"
)

func Run(room *model.Room) {
	for {
		select {
		case client := <-room.Join:
			room.Clients[client] = true
		case client := <-room.Leave:
			delete(room.Clients, client)
			close(client.Send)
		case chat := <-room.Forward:
			for client := range room.Clients {
				client.Send <- chat
			}
		}
	}
}

func NewRoom() *model.Room {
	room := &model.Room{
		Forward: make(chan *domain.Chat),
		Join:    make(chan *model.Client),
		Leave:   make(chan *model.Client),
		Clients: make(map[*model.Client]bool),
	}
	go Run(room)
	return room
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var wsupgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func StartChatRoom(room *model.Room, savedChat *chan model.SavedChat) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id := req.URL.Query().Get("id")

		socket, err := wsupgrader.Upgrade(w, req, nil)
		if err != nil {
			log.Printf("Error while upgrading socket %v", err.Error())
			return
		}

		client := &model.Client{
			Socket:  socket,
			Send:    make(chan *domain.Chat, messageBufferSize),
			Room:    room,
			Channel: id,
			Saved:   savedChat,
		}

		room.Join <- client
		defer func() { room.Leave <- client }()
		go client.Write()
		client.Read()
	}
}
