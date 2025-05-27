package handlers

import (
	"backend/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteRoomHandler(c *gin.Context) {
	log.Println("DeleteRoomHandler called")

	var requestBody struct {
		RoomUUID string `json:"roomUUID"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Printf("Request body: %+v", requestBody)

	queryDeleteRoom := "DELETE FROM rooms WHERE uuid = ?"
	_, err := db.DB.Exec(queryDeleteRoom, requestBody.RoomUUID)
	if err != nil {
		log.Printf("Error deleting room from database: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete room"})
		return
	}

	log.Printf("Room with UUID %s deleted successfully", requestBody.RoomUUID)
	c.JSON(http.StatusOK, gin.H{"message": "Room deleted successfully"})
}
