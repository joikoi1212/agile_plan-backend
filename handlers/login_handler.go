package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	oauth2github "golang.org/x/oauth2/github"
)

var oauth2Config = &oauth2.Config{
	ClientID:     os.Getenv("CLIENT_ID"),
	ClientSecret: os.Getenv("CLIENT_SECRET"),
	Endpoint:     oauth2github.Endpoint,
	RedirectURL:  "https://agileplan-backend-production.up.railway.app/api/v1/callback",
	Scopes:       []string{"read:org", "user"},
}

func LoginHandler(c *gin.Context) {
	oauth2Config.ClientID = strings.TrimSpace(os.Getenv("CLIENT_ID"))
	oauth2Config.ClientSecret = strings.TrimSpace(os.Getenv("CLIENT_SECRET"))
	state := generateState()
	forceLogin, err := c.Request.Cookie("force_github_login")
	var url string
	if err == nil && forceLogin.Value == "1" {
		url = oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOnline, oauth2.SetAuthURLParam("prompt", "login"))
		http.SetCookie(c.Writer, &http.Cookie{
			Name:   "force_github_login",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	} else {
		url = oauth2Config.AuthCodeURL(state)
	}

	fmt.Println("###############", url)
	fmt.Println("CLIENT SECRET 999999999999999", oauth2Config.ClientSecret)
	c.Redirect(http.StatusFound, url)
}

func generateState() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}
