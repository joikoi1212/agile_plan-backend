package handlers

import (
	"backend/db"
	"log"
)

func DeleteRoomHandlerSocket(roomUUID string) error {
	log.Printf("DeleteRoomHandler called for roomUUID: %s", roomUUID)

	queryDeleteRoom := "DELETE FROM rooms WHERE uuid = ?"
	_, err := db.DB.Exec(queryDeleteRoom, roomUUID)
	if err != nil {
		log.Printf("Error deleting room from database: %v", err)
		return err
	}

	log.Printf("Room with UUID %s deleted successfully", roomUUID)
	return nil
}
