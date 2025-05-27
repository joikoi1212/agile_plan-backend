package handlers

import (
	"backend/db"
	"log"
)

type Player struct {
	ID      string
	Name    string
	IsAdmin bool
}

func RemovePlayerByUUID(playerUUID string) ([]Player, error) {

	var roomID string
	err := db.DB.QueryRow("SELECT room_id FROM players WHERE uuid = ?", playerUUID).Scan(&roomID)
	if err != nil {
		log.Printf("Error finding room for player %s: %v", playerUUID, err)
		return nil, err
	}

	_, err = db.DB.Exec("DELETE FROM players WHERE uuid = ?", playerUUID)
	if err != nil {
		log.Printf("Error removing player %s: %v", playerUUID, err)
		return nil, err
	}
	log.Printf("Player with UUID %s removed successfully", playerUUID)

	rows, err := db.DB.Query("SELECT uuid, name, is_admin FROM players WHERE room_id = ?", roomID)
	if err != nil {
		log.Printf("Error retrieving players for room %s: %v", roomID, err)
		return nil, err
	}
	defer rows.Close()

	var players []Player
	for rows.Next() {
		var p Player
		if err := rows.Scan(&p.ID, &p.Name, &p.IsAdmin); err != nil {
			continue
		}
		players = append(players, p)
	}

	return players, nil
}
