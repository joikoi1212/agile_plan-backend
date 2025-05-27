package handlers

import (
	"backend/db"
	"fmt"
	"log"
	"math/rand"
	"time"
)

func CreateRoomHandlerSocket(roomname string, playerName string) (map[string]interface{}, error) {

	log.Println("CreateRoomHandlerSocket called")
	log.Printf("Room name: %s", roomname)
	log.Printf("Player name: %s", playerName)

	roomKey := generateRoomKey(6)

	queryRoom := "INSERT INTO rooms (key_,name) VALUES (?,?)"
	_, err := db.DB.Exec(queryRoom, roomKey, roomname)
	if err != nil {
		log.Printf("Error inserting room into database: %v", err)
		return nil, err
	}

	var roomID string
	queryGetRoomID := "SELECT uuid FROM rooms WHERE name = ?"
	err = db.DB.QueryRow(queryGetRoomID, roomname).Scan(&roomID)
	if err != nil {
		log.Printf("Error retrieving room ID: %v", err)
		return nil, err
	}
	log.Printf("Player name: %s", playerName)
	log.Printf("Room ID: %s", roomID)
	queryPlayer := "INSERT INTO players (room_id, name, is_admin) VALUES (?, ?, ?)"
	_, err = db.DB.Exec(queryPlayer, roomID, playerName, true)
	if err != nil {
		log.Printf("Error assigning player as admin: %v", err)
		return nil, err
	}

	var playerID string
	queryGetPlayerID := "SELECT uuid FROM players WHERE room_id = ? AND name = ?"
	err = db.DB.QueryRow(queryGetPlayerID, roomID, playerName).Scan(&playerID)
	if err != nil {
		log.Printf("Error retrieving player ID: %v", err)
		return nil, err
	}

	response := map[string]interface{}{
		"action": "roomCreated",
		"room": map[string]string{
			"id":   roomID,
			"name": roomname,
			"key":  roomKey,
		},
		"player": map[string]interface{}{
			"id":      playerID,
			"name":    playerName,
			"isAdmin": true,
		},
	}

	return response, nil
}

func generateRoomKey(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	min := int64(100000)
	max := int64(999999)
	return fmt.Sprintf("%06d", r.Int63n(max-min+1)+min)
}
