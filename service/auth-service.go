package service

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
	"github.com/skip2/go-qrcode"
)

type SocketLogin struct {
	Status int
	Socket *websocket.Conn
}

const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var tokens = make(map[string]SocketLogin)

func GenerateToken() (token string) {
	b := make([]byte, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	tokens[string(b)] = SocketLogin{
		Status: 1,
	}
	return string(b)
}

func GenerateQRCode(c *gin.Context) {
	var png []byte
	token := c.Param("token")

	c.Writer.Header().Set("Content-Type", "image/png")
	png, err := qrcode.Encode("http://localhost:9000/api/alpha/auth/token/"+token, qrcode.Medium, 256)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		fmt.Printf("Length is %d bytes long\n", len(png))
	}
	c.Writer.Write(png)
}

func ActivateSocket(c *gin.Context) {
	token := c.Param("token")
	socket, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error while upgrading socket %v", err.Error())
		return
	}

	if tokens[token].Status != 1 {
		defer socket.Close()
		err := socket.WriteMessage(websocket.TextMessage, []byte("Token is not active"))
		if err != nil {
			return
		}
	}

	tokens[token] = SocketLogin{
		Socket: socket,
	}
	for {
		mt, msg, err := socket.ReadMessage()
		if err != nil {
			log.Printf("Error while receiving message %v", err.Error())
			break
		}
		message := "Received " + string(msg)

		if err = socket.WriteMessage(mt, []byte(message)); err != nil {
			log.Printf("Error while sending message %v", err.Error())
			break
		}
	}
}

func ValidateToken(c *gin.Context) {
	token := c.Param("token")
	if tokens[token].Socket != nil {
		socket := tokens[token].Socket

		defer socket.Close()
		err := socket.WriteMessage(websocket.TextMessage, []byte("Token is valid"))
		if err != nil {
			return
		}

		delete(tokens, token)
		log.Println(tokens)
	}
}
