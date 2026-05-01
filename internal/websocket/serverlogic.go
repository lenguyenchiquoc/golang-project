package wsocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ChatMessage struct {
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Room      string `json:"room,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type Client struct {
	Conn     *websocket.Conn
	UserID   string
	Username string
	Room     string
}

type ChatHub struct {
	Clients    map[string]*Client
	Rooms      map[string]map[string]*Client
	Broadcast  chan ChatMessage
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
}

func NewChatHub() *ChatHub {
	return &ChatHub{
		Clients:    make(map[string]*Client),
		Rooms:      make(map[string]map[string]*Client),
		Broadcast:  make(chan ChatMessage, 100),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *ChatHub) Run() {
	for {
		select {
		case client := <-h.Register:

			h.mu.Lock()
			h.Clients[client.Username] = client

			// Add to room
			if h.Rooms[client.Room] == nil {
				h.Rooms[client.Room] = make(map[string]*Client)
			}
			h.Rooms[client.Room][client.Username] = client
			h.mu.Unlock()

			log.Printf("Client joined: %s (room: %s)\n", client.Username, client.Room)

			// Notify others in room
			h.broadcastToRoom(client.Room, ChatMessage{
				Type:      "join",
				Sender:    client.Username,
				Content:   client.Username + " joined the room",
				Room:      client.Room,
				Timestamp: time.Now().Format("15:04:05"),
			}, client.Username)

		case client := <-h.Unregister:
			room := client.Room
			h.mu.Lock()
			delete(h.Clients, client.Username)
			if room, ok := h.Rooms[client.Room]; ok {
				delete(room, client.Username)
				if len(room) == 0 {
					delete(h.Rooms, client.Room)
				}
			}
			h.mu.Unlock()
			client.Conn.Close()

			log.Printf("Client left: %s (room: %s)\n", client.Username, client.Room)
			h.mu.Lock()
			_, roomExists := h.Rooms[room]
			h.mu.Unlock()

			// Notify others
			if roomExists {
				h.broadcastToRoom(client.Room, ChatMessage{
					Type:      "leave",
					Sender:    client.Username,
					Content:   client.Username + " left the room",
					Room:      client.Room,
					Timestamp: time.Now().Format("15:04:05"),
				}, client.Username)
			}

		case msg := <-h.Broadcast:
			h.broadcastToRoom(msg.Room, msg, "")
		}
	}
}

func (h *ChatHub) HandleWS(w http.ResponseWriter, r *http.Request) {
	// Lấy username và room từ query params
	username := r.URL.Query().Get("username")
	room := r.URL.Query().Get("room")

	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}
	if room == "" {
		room = "general"
	}

	// Upgrade HTTP → WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &Client{
		Conn:     conn,
		Username: username,
		Room:     room,
	}

	h.Register <- client

	// Send welcome message
	welcome := ChatMessage{
		Type:      "system",
		Content:   "Welcome to room [" + room + "]!",
		Timestamp: time.Now().Format("15:04:05"),
	}
	data, _ := json.Marshal(welcome)
	conn.WriteMessage(websocket.TextMessage, data)

	// Read messages from client
	defer func() {
		h.Unregister <- client
	}()

	for {
		_, rawMsg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var msg ChatMessage
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			continue
		}

		msg.Sender = username
		msg.Timestamp = time.Now().Format("15:04:05")

		switch msg.Type {
		case "chat":
			msg.Room = client.Room
			h.Broadcast <- msg

		case "dm":
			// Direct message
			h.sendDM(client, msg)

		case "join_room":
			// Switch room
			h.switchRoom(client, msg.Content)

		case "list_rooms":
			h.sendRoomList(client)

		case "list_users":
			h.sendUserList(client)
		}
	}
}

func (h *ChatHub) broadcastToRoom(room string, msg ChatMessage, excludeUsername string) {
	data, _ := json.Marshal(msg)

	h.mu.Lock()
	defer h.mu.Unlock()

	clients, ok := h.Rooms[room]
	if !ok {
		return
	}

	for username, client := range clients {
		if username == excludeUsername {
			continue
		}
		err := client.Conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Printf("Error sending to %s: %v\n", username, err)
		}
	}
}

func (h *ChatHub) sendDM(sender *Client, msg ChatMessage) {
	h.mu.Lock()
	recipient, ok := h.Clients[msg.Recipient]
	h.mu.Unlock()

	if !ok {
		// User not found
		errMsg := ChatMessage{
			Type:      "error",
			Content:   "User " + msg.Recipient + " not found",
			Timestamp: time.Now().Format("15:04:05"),
		}
		data, _ := json.Marshal(errMsg)
		sender.Conn.WriteMessage(websocket.TextMessage, data)
		return
	}

	msg.Type = "dm"
	data, _ := json.Marshal(msg)
	recipient.Conn.WriteMessage(websocket.TextMessage, data)
}

func (h *ChatHub) switchRoom(client *Client, newRoom string) {
	oldRoom := client.Room

	h.mu.Lock()
	// Remove from old room
	if room, ok := h.Rooms[oldRoom]; ok {
		delete(room, client.Username)
		if len(room) == 0 {
			delete(h.Rooms, oldRoom)
		}
	}
	// Add to new room
	if h.Rooms[newRoom] == nil {
		h.Rooms[newRoom] = make(map[string]*Client)
	}
	h.Rooms[newRoom][client.Username] = client
	client.Room = newRoom
	h.mu.Unlock()

	// Notify new room
	h.broadcastToRoom(newRoom, ChatMessage{
		Type:      "join",
		Sender:    client.Username,
		Content:   client.Username + " joined the room",
		Room:      newRoom,
		Timestamp: time.Now().Format("15:04:05"),
	}, client.Username)

	// Confirm to client
	confirm := ChatMessage{
		Type:      "system",
		Content:   "Switched to room [" + newRoom + "]",
		Timestamp: time.Now().Format("15:04:05"),
	}
	data, _ := json.Marshal(confirm)
	client.Conn.WriteMessage(websocket.TextMessage, data)
}

func (h *ChatHub) sendRoomList(client *Client) {
	h.mu.Lock()
	rooms := []string{}
	for room := range h.Rooms {
		rooms = append(rooms, room)
	}
	h.mu.Unlock()

	msg := ChatMessage{
		Type:      "list_rooms",
		Content:   join(rooms, ", "),
		Timestamp: time.Now().Format("15:04:05"),
	}
	data, _ := json.Marshal(msg)
	client.Conn.WriteMessage(websocket.TextMessage, data)
}

func (h *ChatHub) sendUserList(client *Client) {
	h.mu.Lock()
	users := []string{}
	if room, ok := h.Rooms[client.Room]; ok {
		for username := range room {
			users = append(users, username)
		}
	}
	h.mu.Unlock()

	msg := ChatMessage{
		Type:      "user_list",
		Content:   join(users, ", "),
		Timestamp: time.Now().Format("15:04:05"),
	}
	data, _ := json.Marshal(msg)
	client.Conn.WriteMessage(websocket.TextMessage, data)
}

func join(items []string, sep string) string {
	result := ""
	for i, item := range items {
		if i > 0 {
			result += sep
		}
		result += item
	}
	return result
}
