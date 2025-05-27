package websocket

import (
	"backend/handlers"
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	send       chan []byte
	mu         sync.Mutex
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		send:       make(chan []byte, 1024),
	}
}

func (c *Client) readMessages() {
	log.Println("Client readMessages started")
	defer func() {
		c.manager.removeClient(c)
		c.connection.Close()
	}()

	for {
		_, message, err := c.connection.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Raw message received: %s", message)
		var payload map[string]interface{}
		if err := json.Unmarshal(message, &payload); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		log.Printf("Parsed payload: %+v", payload)
		action, ok := payload["action"].(string)
		if !ok {
			log.Printf("Missing or invalid action in message")
			continue
		}

		switch action {

		case "createRoom":
			log.Println("Create room action received")
			roomname, _ := payload["roomname"].(string)
			playerName, _ := payload["playerName"].(string)
			log.Printf("Player name: %s", playerName)

			response, err := handlers.CreateRoomHandlerSocket(roomname, playerName)
			if err != nil {
				c.SendMessage([]byte(`{"action": "error", "message": "Failed to create room"}`))
				continue
			}

			roomUUID := response["room"].(map[string]string)["id"]
			c.manager.addClientToRoom(c, roomUUID)
			log.Printf("Admin added to room %s", roomUUID)
			responseBytes, _ := json.Marshal(response)
			c.SendMessage(responseBytes)
			c.manager.broadcastToRoom(roomUUID, responseBytes)

		case "joinRoom":
			log.Println("Join room action received")
			roomKey, _ := payload["roomKey"].(string)

			player, ok := payload["player"].(map[string]interface{})
			if !ok {
				log.Println("Invalid player object in payload")
				c.SendMessage([]byte(`{"action": "error", "message": "Invalid player object"}`))
				continue
			}

			playerName, _ := player["name"].(string)
			if playerName == "" {
				log.Println("Player name is empty or missing")
				c.SendMessage([]byte(`{"action": "error", "message": "Player name cannot be empty"}`))
				continue
			}

			log.Printf("RoomKey: %s, PlayerName: %s", roomKey, playerName)

			c.manager.addClientToRoom(c, roomKey)
			log.Printf("Client joining room %s", roomKey)

			response, err := handlers.JoinRoomHandlerSocket(roomKey, playerName)
			if err != nil {
				c.SendMessage([]byte(`{"action": "error", "message": "Failed to join room"}`))
				continue
			}

			roomUUID := response["room"].(map[string]string)["id"]
			c.manager.addClientToRoom(c, roomUUID)
			log.Printf("Client joining room %s", roomUUID)
			responseBytes, _ := json.Marshal(response)
			c.SendMessage(responseBytes)
			c.manager.broadcastToRoom(roomUUID, responseBytes)

		case "deleteRoom":
			log.Println("Delete room action received")
			roomUUID, _ := payload["roomUUID"].(string)
			if roomUUID == "" {
				log.Println("RoomUUID is empty or missing")
				c.SendMessage([]byte(`{"action": "error", "message": "RoomUUID cannot be empty"}`))
				break
			}

			err := handlers.DeleteRoomHandlerSocket(roomUUID)
			if err != nil {
				log.Printf("Error deleting room: %v", err)
				c.SendMessage([]byte(`{"action": "error", "message": "Failed to delete room"}`))
				break
			}

			updateMsg := map[string]interface{}{
				"action":   "roomClosed",
				"roomUUID": roomUUID,
			}
			updateBytes, _ := json.Marshal(updateMsg)
			c.manager.broadcastToRoom(roomUUID, updateBytes)

		case "leaveRoom":
			log.Println("Leave room action received")
			roomUUID, _ := payload["roomUUID"].(string)
			playerUUID, _ := payload["playerUUID"].(string)
			if roomUUID == "" || playerUUID == "" {
				log.Println("Missing roomUUID or playerUUID")
				c.SendMessage([]byte(`{"action": "error", "message": "Missing roomUUID or playerUUID"}`))
				break
			}

			players, err := handlers.RemovePlayerByUUID(playerUUID)
			if err != nil {
				log.Printf("Error removing player: %v", err)
				c.SendMessage([]byte(`{"action": "error", "message": "Failed to leave room"}`))
				break
			}

			playerList := []map[string]interface{}{}
			for _, p := range players {
				playerList = append(playerList, map[string]interface{}{
					"id":      p.ID,
					"name":    p.Name,
					"isAdmin": p.IsAdmin,
				})
			}
			updateMsg := map[string]interface{}{
				"action":   "playerLeft",
				"roomUUID": roomUUID,
				"players":  playerList,
				"leaverId": playerUUID,
			}
			updateBytes, _ := json.Marshal(updateMsg)
			c.manager.broadcastToRoom(roomUUID, updateBytes)

			c.manager.removeClientFromRoom(c, roomUUID)

		case "assign_ticket":
			log.Println("Assign ticket action received")
			roomUUID, _ := payload["roomUUID"].(string)
			ticketKey, _ := payload["ticketKey"].(string)
			if roomUUID == "" || ticketKey == "" {
				log.Println("Missing roomUUID or ticketKey")
				c.SendMessage([]byte(`{"action": "error", "message": "Missing roomUUID or ticketKey"}`))
				break
			}

			err := handlers.AssignTicketToRoom(roomUUID, ticketKey)
			if err != nil {
				log.Printf("Error assigning ticket: %v", err)
				c.SendMessage([]byte(`{"action": "error", "message": "Failed to assign ticket"}`))
				break
			}
			updateMsg := map[string]interface{}{
				"action":    "ticketAssigned",
				"roomKey":   roomUUID,
				"ticketKey": ticketKey,
			}
			updateBytes, _ := json.Marshal(updateMsg)
			c.manager.broadcastToRoom(roomUUID, updateBytes)

		case "start_vote":
			roomUUID, _ := payload["roomUUID"].(string)
			c.manager.votes[roomUUID] = make(map[string]int)
			c.manager.revealed[roomUUID] = false
			updateMsg := map[string]interface{}{
				"action":   "voteStarted",
				"roomUUID": roomUUID,
			}
			updateBytes, _ := json.Marshal(updateMsg)
			c.manager.broadcastToRoom(roomUUID, updateBytes)

		case "vote":
			roomUUID, _ := payload["roomUUID"].(string)
			playerId, _ := payload["playerId"].(string)
			value, _ := payload["value"].(float64)
			if _, ok := c.manager.votes[roomUUID]; !ok {
				c.manager.votes[roomUUID] = make(map[string]int)
			}
			c.manager.votes[roomUUID][playerId] = int(value)
			voted := map[string]bool{}
			for pid := range c.manager.votes[roomUUID] {
				voted[pid] = true
			}
			updateMsg := map[string]interface{}{
				"action": "voteUpdate",
				"votes":  voted,
			}
			updateBytes, _ := json.Marshal(updateMsg)
			c.manager.broadcastToRoom(roomUUID, updateBytes)

		case "revealVotes":
			roomUUID, _ := payload["roomUUID"].(string)
			c.manager.revealed[roomUUID] = true
			votes := c.manager.votes[roomUUID]
			updateMsg := map[string]interface{}{
				"action": "votesRevealed",
				"votes":  votes,
			}
			updateBytes, _ := json.Marshal(updateMsg)
			c.manager.broadcastToRoom(roomUUID, updateBytes)

		case "returnToRoom":
			roomUUID, _ := payload["roomUUID"].(string)
			updateMsg := map[string]interface{}{
				"action":   "returnToRoom",
				"roomUUID": roomUUID,
			}
			updateBytes, _ := json.Marshal(updateMsg)
			c.manager.broadcastToRoom(roomUUID, updateBytes)

		case "fibApproxChanged":
			roomUUID, _ := payload["roomUUID"].(string)
			direction, _ := payload["direction"].(string)
			updateMsg := map[string]interface{}{
				"action":    "fibApproxChanged",
				"roomUUID":  roomUUID,
				"direction": direction,
			}
			updateBytes, _ := json.Marshal(updateMsg)
			c.manager.broadcastToRoom(roomUUID, updateBytes)

		case "addRoomHistory":
			roomUUID, _ := payload["roomUUID"].(string)
			ticketName, _ := payload["ticketName"].(string)
			date, _ := payload["date"].(string)
			media := payload["media"]

			c.manager.Lock()
			c.manager.roomHistory[roomUUID] = append(c.manager.roomHistory[roomUUID], RoomHistoryItem{
				TicketName: ticketName,
				Date:       date,
				Media:      media,
			})
			history := c.manager.roomHistory[roomUUID]
			c.manager.Unlock()

			response := map[string]interface{}{
				"action":   "roomHistoryUpdated",
				"roomUUID": roomUUID,
				"history":  history,
			}
			responseBytes, _ := json.Marshal(response)
			c.manager.broadcastToRoom(roomUUID, responseBytes)

		case "getRoomHistory":
			roomUUID, _ := payload["roomUUID"].(string)
			c.manager.RLock()
			history := c.manager.roomHistory[roomUUID]
			c.manager.RUnlock()

			response := map[string]interface{}{
				"action":   "roomHistoryUpdated",
				"roomUUID": roomUUID,
				"history":  history,
			}
			responseBytes, _ := json.Marshal(response)
			c.SendMessage(responseBytes)

		default:
			log.Printf("Unknown action: %s", action)
		}
	}
}

func (c *Client) writeMessages() {
	defer c.connection.Close()

	for message := range c.send {
		c.mu.Lock()
		err := c.connection.WriteMessage(websocket.TextMessage, message)
		c.mu.Unlock()

		if err != nil {
			log.Printf("Error writing message: %v", err)
			break
		}
	}
}

func (c *Client) SendMessage(message []byte) {
	select {
	case c.send <- message:
		log.Printf("Message sent to client: %s", string(message))
	default:
		log.Println("Send buffer full, closing connection")
		c.manager.removeClient(c)
		c.connection.Close()
	}
}
