package controller

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ivohutasoit/alira/util"
)

func IndexPageHandler(c *gin.Context) {
	url, err := util.GenerateUrl(c.Request.TLS, c.Request.Host, c.Request.URL.Path, true)
	if err != nil {
		fmt.Println(err)
		return
	}
	redirect := fmt.Sprintf("%s?redirect=%s", os.Getenv("LOGOUT_URL"), url)
	c.HTML(http.StatusOK, "home.tmpl.html", gin.H{
		"logout_url": redirect,
	})
}
