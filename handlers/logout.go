package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   "force_github_login",
		Value:  "1",
		Path:   "/",
		MaxAge: 20,
	})
	c.Redirect(http.StatusFound, "https://github.com/logout")
}
