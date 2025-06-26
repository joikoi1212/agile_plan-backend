package routes

import (
	"backend/handlers"
	"backend/websocket"
	"log"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, manager *websocket.Manager) {
	// WebSocket route is now handled in main.go before middleware
	
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
