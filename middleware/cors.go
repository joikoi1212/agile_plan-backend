package middleware

import (
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	hostname := os.Getenv("FRONTEND_HOSTNAME")
	port := os.Getenv("FRONTEND_PORT")
	fmt.Println("CORS hostname:", hostname)
	fmt.Println("CORS port:", port)
	allowedOrigin := fmt.Sprintf("http://%s:%s", hostname, port)
	fmt.Printf("CORS Middleware: Allowing origin %s\n", allowedOrigin)
	return cors.New(cors.Config{
		AllowOrigins:     []string{"https://agile-plan-backend.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	})
}
