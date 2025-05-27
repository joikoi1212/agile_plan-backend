package handlers

import (
	"backend/db"
	"log"
)

func JoinRoomHandlerSocket(roomKey string, playerName string) (map[string]interface{}, error) {
	var roomUUID, roomName string

	err := db.DB.QueryRow("SELECT uuid, name FROM rooms WHERE key_ = ?", roomKey).Scan(&roomUUID, &roomName)
	if err != nil {
		log.Printf("Room not found: %v", err)
		return nil, err
	}

	_, err = db.DB.Exec("INSERT INTO players (room_id, name, is_admin) VALUES (?, ?, ?)", roomUUID, playerName, false)
	if err != nil {
		log.Printf("Error adding player: %v", err)
		return nil, err
	}

	var playerUUID string
	err = db.DB.QueryRow("SELECT uuid FROM players WHERE room_id = ? AND name = ? ORDER BY id DESC LIMIT 1", roomUUID, playerName).Scan(&playerUUID)
	if err != nil {
		log.Printf("Error retrieving player UUID: %v", err)
		return nil, err
	}

	rows, err := db.DB.Query("SELECT uuid, name, is_admin FROM players WHERE room_id = ?", roomUUID)
	if err != nil {
		log.Printf("Error retrieving players: %v", err)
		return nil, err
	}
	defer rows.Close()

	players := []map[string]interface{}{}
	for rows.Next() {
		var playerID, playerName string
		var isAdmin bool
		if err := rows.Scan(&playerID, &playerName, &isAdmin); err != nil {
			log.Printf("Error scanning player: %v", err)
			return nil, err
		}
		players = append(players, map[string]interface{}{
			"id":      playerID,
			"name":    playerName,
			"isAdmin": isAdmin,
		})
	}

	response := map[string]interface{}{
		"action": "roomJoined",
		"room": map[string]string{
			"id":   roomUUID,
			"key":  roomKey,
			"name": roomName,
		},
		"player": map[string]interface{}{
			"id":      playerUUID,
			"name":    playerName,
			"isAdmin": false,
		},
		"players": players,
	}
	log.Printf("Response SOCKETJOINHANDLER: %+v", response)

	return response, nil
}
