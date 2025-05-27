package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type JoinRoomRequest struct {
	RoomKey    string `json:"roomKey"`
	PlayerName string `json:"playerName"`
}

type JoinRoomResponse struct {
	Room struct {
		ID  string `json:"id"`
		Key string `json:"key"`
	} `json:"room"`
	Player struct {
		Name    string `json:"name"`
		IsAdmin bool   `json:"isAdmin"`
	} `json:"player"`
}

type RoomHistoryItem struct {
	TicketName string      `json:"ticketName"`
	Date       string      `json:"date"`
	Media      interface{} `json:"media"`
}

type Manager struct {
	clients     map[*Client]bool
	rooms       map[string]map[*Client]bool
	votes       map[string]map[string]int
	revealed    map[string]bool
	roomHistory map[string][]RoomHistoryItem
	sync.RWMutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {

		return true
	},
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := NewClient(conn, m)
	m.addClient(client)

	go client.readMessages()
	go client.writeMessages()
}

func NewManager() *Manager {
	return &Manager{
		clients:     make(map[*Client]bool),
		rooms:       make(map[string]map[*Client]bool),
		votes:       make(map[string]map[string]int),
		revealed:    make(map[string]bool),
		roomHistory: make(map[string][]RoomHistoryItem),
	}
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true
	log.Printf("Client added. Total clients: %d", len(m.clients))
}

func sendError(conn *websocket.Conn, message string) {
	errorResponse := map[string]string{
		"action":  "error",
		"message": message,
	}
	responseBytes, _ := json.Marshal(errorResponse)
	conn.WriteMessage(websocket.TextMessage, responseBytes)
}

func (m *Manager) addClientToRoom(client *Client, roomUUID string) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.rooms[roomUUID]; !ok {
		m.rooms[roomUUID] = make(map[*Client]bool)
	}

	if _, exists := m.rooms[roomUUID][client]; exists {
		log.Printf("Client already in room %s. Skipping addition.", roomUUID)
		return
	}

	m.rooms[roomUUID][client] = true
	log.Printf("Client %p added to room %s", client, roomUUID)
	log.Printf("Client added to room %s. Total clients in room: %d", roomUUID, len(m.rooms[roomUUID]))
}

func (m *Manager) removeClientFromRoom(client *Client, roomUUID string) {
	m.Lock()
	defer m.Unlock()

	if clients, ok := m.rooms[roomUUID]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(m.rooms, roomUUID)
		}
	}
}

func (m *Manager) broadcastToRoom(roomUUID string, message []byte) {
	m.RLock()
	defer m.RUnlock()

	if clients, ok := m.rooms[roomUUID]; ok {
		log.Printf("Broadcasting message to room %s with %d clients", roomUUID, len(clients))

		for client := range clients {
			go func(c *Client) {
				c.SendMessage(message)
			}(client)
		}
	} else {
		log.Printf("No clients found in room %s", roomUUID)
	}
}
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	for roomKey, clients := range m.rooms {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			if len(clients) == 0 {
				delete(m.rooms, roomKey)
			}
		}
	}

	delete(m.clients, client)
}
