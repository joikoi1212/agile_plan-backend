package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("userID")

		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired. Please log in again."})
			c.Abort()
			return
		}

		c.Next()
	}
}
