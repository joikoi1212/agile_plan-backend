package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v70/github"
	"golang.org/x/oauth2"
)

func CallbackHandler(c *gin.Context) {
	code := c.Query("code")
	fmt.Println("CLIENT SECRET $$$$$$$$$$$$", oauth2Config.ClientSecret)
	token, err := oauth2Config.Exchange(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		fmt.Println(err)
		fmt.Println("###############################")
		fmt.Println("###################################       ", code)
		return
	}

	githubUser, err := getGithubUser(token)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Failed to get user info"})
		return
	}

	fmt.Println("User ID:", *githubUser.Login)
	fmt.Println("User Name:", *githubUser.Name)
	fmt.Println("Avatar URL:", *githubUser.AvatarURL)
	session := sessions.Default(c)
	session.Set("id", *githubUser.Login)
	session.Set("username", *githubUser.Name)
	session.Set("avatar", *githubUser.AvatarURL)
	session.Save()
	fmt.Println("Session username:", session.Get("username"))
	fmt.Println("Session avatar:", session.Get("avatar"))
	fmt.Println("Session ID:", session.Get("id"))
	c.Redirect(http.StatusFound, "https://your-project-name.vercel.app/dashboard")

}

func getGithubUser(token *oauth2.Token) (*github.User, error) {
	ctx := context.Background()
	client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(token)))

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, err
	}
	return user, nil
}
