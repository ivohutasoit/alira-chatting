package service

import (
	"log"

	"github.com/ivohutasoit/alira-chatting/model"
)

func NewSavedChat() *chan model.SavedChat {
	chat := make(chan model.SavedChat, 256)
	go SaveChat(&chat)
	return &chat
}

func SaveChat(chat *chan model.SavedChat) {
	for {
		_, ok := <-*chat
		if !ok {
			log.Print("Error when receiving chat to save")
			return
		}
	}
}
