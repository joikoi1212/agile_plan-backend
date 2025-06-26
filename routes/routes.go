package routes

import (
	"backend/handlers"
	"backend/websocket"
	"log"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, manager *websocket.Manager) {
	// Register WebSocket route FIRST (before other routes that might apply middleware)
	r.GET("/ws", func(c *gin.Context) {
		// Log all WebSocket connection attempts
		log.Printf("WebSocket connection attempt from %s, Origin: %s", c.ClientIP(), c.GetHeader("Origin"))
		
		// Add CORS headers for WebSocket
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		
		// Initialize session if it doesn't exist (but don't fail if it does)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Session recovery in WebSocket: %v (continuing)", r)
			}
		}()
		
		manager.ServeWS(c.Writer, c.Request)
	})

	// Test endpoint to check if route is reachable
	r.GET("/ws-test", func(c *gin.Context) {
		log.Printf("WebSocket test endpoint hit from %s", c.ClientIP())
		c.JSON(200, gin.H{"message": "WebSocket endpoint is reachable", "ready": true})
	})

	group := r.Group("/api/v1")
	RegisterAuthRoutes(group)
}

func RegisterAuthRoutes(rg *gin.RouterGroup) {
	rg.GET("/session", handlers.SessionHandler)
	rg.GET("/login", handlers.LoginHandler)
	rg.GET("/callback", handlers.CallbackHandler)
	rg.GET("/tickets", handlers.TicketsHandler)
	rg.POST("/createroom", handlers.CreateRoomHandler)
	rg.POST("/deleteroom", handlers.DeleteRoomHandler)
	rg.POST("/joinroom", handlers.JoinRoomHandler)
	rg.GET("/logout", handlers.LogoutHandler)
}
