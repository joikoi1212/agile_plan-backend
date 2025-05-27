package handlers

import (
	"backend/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JoinRoomHandler(c *gin.Context) {
	log.Println("JoinRoomHandler called")

	var requestBody struct {
		RoomKey    string `json:"roomkey"`
		PlayerName string `json:"playerName"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("Error parsing request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var roomUUID string
	err := db.DB.QueryRow("SELECT uuid FROM rooms WHERE key_ = ?", requestBody.RoomKey).Scan(&roomUUID)
	if err != nil {
		log.Printf("Room not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	_, err = db.DB.Exec("INSERT INTO players (room_id, name, is_admin) VALUES (?, ?, ?)", roomUUID, requestBody.PlayerName, false)
	if err != nil {
		log.Printf("Error adding player: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"room": gin.H{
			"id":  roomUUID,
			"key": requestBody.RoomKey,
		},
		"player": gin.H{
			"name":    requestBody.PlayerName,
			"isAdmin": false,
		},
	})
}
