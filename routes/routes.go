package routes

import (
	"backend/handlers"
	"backend/websocket"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, manager *websocket.Manager) {
	group := r.Group("/api/v1")
	RegisterAuthRoutes(group)

	r.GET("/ws", func(c *gin.Context) {
		manager.ServeWS(c.Writer, c.Request)
	})
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
