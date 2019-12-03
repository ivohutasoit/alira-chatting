package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/ivohutasoit/alira-chatting/service"
	"github.com/ivohutasoit/alira/middleware"
	"github.com/ivohutasoit/alira/util"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		fmt.Println("$PORT must be set")
		port = "9001"
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(sessions.Sessions("ALIRASESSIOn", util.Store))
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		/*c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"token": service.GenerateToken(),
		})*/
		http.ServeFile(c.Writer, c.Request, "templates/home.html")
	})

	api := router.Group("/api/alpha")
	{
		room := service.CreateChatRoom()
		chat := service.CreateSavedChat()
		api.Use(middleware.AuthenticationRequired())

		api.GET("/channel/:id", func(c *gin.Context) {
			service.StartChatRoom(room, chat)
		})
	}

	router.Run(":" + port)
}
