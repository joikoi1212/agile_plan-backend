package handlers

import (
	"backend/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRoomHandler(c *gin.Context) {
	log.Println("CreateRoomHandler called")

	var requestBody struct {
		Roomname   string `json:"roomname"`
		PlayerName string `json:"playerName"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Printf("Request body: %+v", requestBody)

	roomKey := generateRoomKey(6)
	log.Printf("Generated room key: %s", roomKey)
	queryRoom := "INSERT INTO rooms (key_, name) VALUES (?, ?)"
	_, err := db.DB.Exec(queryRoom, roomKey, requestBody.Roomname)
	if err != nil {
		log.Printf("Error inserting room into database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
		return
	}

	var roomID string
	queryGetRoomID := "SELECT id FROM rooms WHERE key_ = ?"
	err = db.DB.QueryRow(queryGetRoomID, roomKey).Scan(&roomID)
	if err != nil {
		log.Printf("Error retrieving room ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve room ID"})
		return
	}

	var playerID string
	queryPlayer := "INSERT INTO players (room_id, name, is_admin) VALUES (?, ?, ?) RETURNING id"
	err = db.DB.QueryRow(queryPlayer, roomID, requestBody.PlayerName, true).Scan(&playerID)
	if err != nil {
		log.Printf("Error assigning player as admin: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign player as admin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"room": gin.H{
			"id":   roomID,
			"key":  roomKey,
			"name": requestBody.Roomname,
		},
		"player": gin.H{
			"id":      playerID,
			"name":    requestBody.PlayerName,
			"isAdmin": true,
		},
	})

	log.Printf("Room created: %+v", gin.H{
		"room": gin.H{
			"id":   roomID,
			"key":  roomKey,
			"name": requestBody.Roomname,
		},
		"player": gin.H{
			"id":      playerID,
			"name":    requestBody.PlayerName,
			"isAdmin": true,
		},
	})
}
