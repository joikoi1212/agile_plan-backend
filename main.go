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
	
	// Create WebSocket manager
	manager := websocket.NewManager()
	
	// Register WebSocket routes BEFORE session middleware
	r.GET("/ws", func(c *gin.Context) {
		log.Printf("WebSocket connection attempt from %s, Origin: %s", c.ClientIP(), c.GetHeader("Origin"))
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		manager.ServeWS(c.Writer, c.Request)
	})
	
	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Server is running"})
	})
	
	// Now apply middleware for other routes
	store := cookie.NewStore([]byte("secret-key"))
	r.Use(sessions.Sessions("my-session", store))
	r.Use(middleware.CORSMiddleware())
	
	// Register API routes (these will use middleware)
	routes.RegisterRoutes(r, manager)

	log.Printf("Server running on port %s\n", port)
	r.Run("0.0.0.0:" + port)
}
