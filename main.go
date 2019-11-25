package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/heroku/x/hmetrics/onload"

	"github.com/ivohutasoit/alira-chatting/model"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		//port = "9000"
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/ws", func(c *gin.Context) {
		handleConnection(c.Writer, c.Request)
	})

	go handleMessages()

	router.Run(":" + port)
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan model.Message)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		//fmt.Println("Failed to set web socket upgrade: %+v", err)
		//return
	}

	defer conn.Close()

	// Register new client
	clients[conn] = true

	for {
		var msg model.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error Read: %v", err)
			delete(clients, conn)
			break
		}
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("Error Write: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
