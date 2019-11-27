package main

import (
	"net/http"
	"os"

	"github.com/ivohutasoit/alira-chatting/service"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "9000"
		//log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"token": service.GenerateToken(),
		})
		//http.ServeFile(c.Writer, c.Request, "templates/home.tmpl.html")
	})
	router.GET("/api/alpha/qrcode/generate/:token", service.GenerateQRCode)
	router.GET("/api/alpha/auth/socket/:token", service.ActivateSocket)
	router.GET("/api/alpha/auth/token/:token", service.ValidateToken)

	room := service.CreateChatRoom()
	chat := service.CreateSavedChat()

	router.GET("/channel/:id", func(c *gin.Context) {
		service.StartChatRoom(room, chat)
	})

	router.Run(":" + port)
}
