package main

import (
	"backend/db"
	"backend/middleware"
	"backend/routes"
	"backend/websocket"
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	db.InitDB()
	store := cookie.NewStore([]byte("secret-key"))
	r.Use(sessions.Sessions("my-session", store))
	r.Use(middleware.CORSMiddleware())
	manager := websocket.NewManager()
	routes.RegisterRoutes(r, manager)

	log.Println("Server running on http://localhost:8088")
	r.Run(":8088")
}
