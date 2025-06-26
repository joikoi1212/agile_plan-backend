package handlers

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SessionHandler(c *gin.Context) {
	log.Println("SessionHandler called")
	session := sessions.Default(c)
	username := session.Get("username")
	avatar := session.Get("avatar")

	if username == nil || avatar == nil {
		log.Println("No session found - treating as guest")
		c.JSON(http.StatusOK, gin.H{
			"guest": true,
			"authenticated": false,
			"message": "Guest mode",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"avatar":   avatar,
		"authenticated": true,
		"guest": false,
	})
}
