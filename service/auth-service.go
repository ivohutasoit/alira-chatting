package service

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ivohutasoit/alira/util"

	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
	"github.com/skip2/go-qrcode"
)

type SocketLogin struct {
	Status int
	Socket *websocket.Conn
}

const (
	letters   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	secretKey = "time2SleepW3LLOK"
)

var tokens = make(map[string]SocketLogin)

func GenerateToken() (token string) {
	b := make([]byte, 16)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	encrypted, err := util.Encrypt([]byte(secretKey), b)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	log.Println(string(b))
	tokens[string(b)] = SocketLogin{
		Status: 1,
	}
	return encrypted
}

func GenerateQRCode(c *gin.Context) {
	var png []byte
	token := c.Param("token")

	c.Writer.Header().Set("Content-Type", "image/png")
	png, err := qrcode.Encode(token, qrcode.Medium, 256)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		fmt.Printf("Length is %d bytes long\n", len(png))
	}
	c.Writer.Write(png)
}

func ActivateSocket(c *gin.Context) {
	token := c.Param("token")
	log.Println(token)
	decrypted, err := util.Decrypt([]byte(secretKey), []byte(token))
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
	socket, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error while upgrading socket %v", err.Error())
		return
	}

	if tokens[decrypted].Status != 1 {
		defer socket.Close()
		err := socket.WriteMessage(websocket.TextMessage, []byte("Token is not active"))
		if err != nil {
			return
		}
	}
	log.Println(decrypted)
	tokens[decrypted] = SocketLogin{
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
