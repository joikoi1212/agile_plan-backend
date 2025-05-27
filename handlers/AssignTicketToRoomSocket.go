package handlers

import (
	"backend/db"
	"log"
)

func AssignTicketToRoom(roomUUID string, ticketKey string) error {
	_, err := db.DB.Exec("UPDATE rooms SET current_ticket = ? WHERE uuid = ?", ticketKey, roomUUID)
	if err != nil {
		log.Printf("Error assigning ticket %s to room %s: %v", ticketKey, roomUUID, err)
		return err
	}
	log.Printf("Assigned ticket %s to room %s", ticketKey, roomUUID)
	return nil
}
