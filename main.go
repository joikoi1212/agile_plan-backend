package main

import (
	"backend/db"
	"backend/middleware"
	"backend/routes"
	"backend/websocket"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r := gin.Default()
	db.InitDB()
	store := cookie.NewStore([]byte("secret-key"))
	r.Use(sessions.Sessions("my-session", store))
	r.Use(middleware.CORSMiddleware())
	manager := websocket.NewManager()
	routes.RegisterRoutes(r, manager)

	log.Printf("Server running on port %s\n", port)
	r.Run(":" + port)
}
