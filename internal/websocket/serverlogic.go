package wsocket

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ChatMessage struct {
	ID        string `json:"id,omitempty"`
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
	send     chan []byte
	hub      *ChatHub
}

type ChatHub struct {
	Clients          map[string]*Client
	Rooms            map[string]map[string]*Client
	UsernameToUserID map[string]string
	Broadcast        chan ChatMessage
	Register         chan *Client
	Unregister       chan *Client
	mu               sync.RWMutex
}

func NewChatHub() *ChatHub {
	return &ChatHub{
		Clients:          make(map[string]*Client),
		Rooms:            make(map[string]map[string]*Client),
		UsernameToUserID: make(map[string]string),
		Broadcast:        make(chan ChatMessage, 512),
		Register:         make(chan *Client, 128),
		Unregister:       make(chan *Client, 128),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = 45 * time.Second
	maxMessageSize = 512 * 1024
)

// ================= CLIENT =================

func (c *Client) readPump() {
	defer func() {
		c.hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, raw, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Read error from %s: %v", c.Username, err)
			}
			return
		}

		var msg ChatMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}

		msg.Sender = c.Username
		msg.Timestamp = time.Now().Format("15:04:05")

		switch msg.Type {
		case "chat":
			msg.Room = c.Room
			c.hub.Broadcast <- msg

		case "dm":
			c.hub.sendDM(c, msg)

		case "join_room":
			c.hub.switchRoom(c, normalizeRoom(msg.Content))

		case "list_rooms":
			c.hub.sendRoomList(c)

		case "list_users":
			c.hub.sendUserList(c)

		case "list_all_users":
			c.hub.sendAllUsers(c)
		}

	}
}

// func (h *ChatHub) sendAllUsers(c *Client) {
// 	h.mu.RLock()
// 	users := make([]string, 0, len(h.Clients))

// 	for _, client := range h.Clients {
// 		if client.UserID != c.UserID {
// 			users = append(users, client.Username)
// 		}
// 	}
// 	h.mu.RUnlock()

// 	c.send <- mustJSON(ChatMessage{
// 		Type:      "list_all_users",
// 		Content:   strings.Join(users, ", "),
// 		Timestamp: time.Now().Format("15:04:05"),
// 	})
// }

func (h *ChatHub) sendAllUsers(c *Client) {
	h.mu.RLock()
	users := make([]string, 0, len(h.Clients))
	for _, client := range h.Clients {
		if client != nil && client.UserID != c.UserID {
			users = append(users, client.Username)
		}
	}
	h.mu.RUnlock()

	log.Printf("📤 Sending %d online users to %s", len(users), c.Username) // debug

	safeSend(c, mustJSON(ChatMessage{
		Type:      "list_all_users",
		Content:   strings.Join(users, ", "),
		Timestamp: time.Now().Format("15:04:05"),
	}))
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ================= HUB =================

func (h *ChatHub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.register(c)
		case c := <-h.Unregister:
			h.unregister(c)
		case msg := <-h.Broadcast:
			h.broadcast(msg)
		}
	}
}

// func (h *ChatHub) register(c *Client) {
// 	h.mu.Lock()
// 	defer h.mu.Unlock()

// 	// Close old session if exists
// 	if old, exists := h.Clients[c.UserID]; exists {
// 		close(old.send)
// 	}

// 	c.Room = normalizeRoom(c.Room)
// 	h.Clients[c.UserID] = c
// 	h.UsernameToUserID[c.Username] = c.UserID

// 	if h.Rooms[c.Room] == nil {
// 		h.Rooms[c.Room] = make(map[string]*Client)
// 	}
// 	h.Rooms[c.Room][c.UserID] = c

// 	log.Printf("✅ %s joined room [%s]", c.Username, c.Room)

// 	// Welcome
// 	welcome := ChatMessage{
// 		Type:      "system",
// 		Content:   "Welcome to room [" + c.Room + "]!",
// 		Timestamp: time.Now().Format("15:04:05"),
// 	}
// 	select {
// 	case c.send <- mustJSON(welcome):
// 	default:
// 	}
// }

func (h *ChatHub) register(c *Client) {
    h.mu.Lock()
    defer h.mu.Unlock()

    c.Room = normalizeRoom(c.Room)
    isPrivate := strings.HasPrefix(c.Room, "private_")

    if old, exists := h.Clients[c.UserID]; exists {
        oldIsPrivate := strings.HasPrefix(old.Room, "private_")

        if isPrivate && oldIsPrivate {
            // Cả 2 đều private → close old (reconnect case)
            if old.send != nil {
                close(old.send)
                old.send = nil
            }
        } else if !isPrivate && !oldIsPrivate {
            // Cả 2 đều public → close old
            if old.send != nil {
                close(old.send)
                old.send = nil
            }
        }
        // Một private một public → giữ cả 2, chỉ update Clients nếu là private
        // Public room connection KHÔNG ghi đè Clients[userID]
    }

    if isPrivate {
        // Private connection là "chính" → lưu vào Clients
        h.Clients[c.UserID] = c
        h.UsernameToUserID[c.Username] = c.UserID
    } else {
        // Public connection → chỉ thêm vào Rooms, không ghi đè Clients
        // nếu đã có private connection
        if _, hasExisting := h.Clients[c.UserID]; !hasExisting {
            h.Clients[c.UserID] = c
            h.UsernameToUserID[c.Username] = c.UserID
        }
        if h.Rooms[c.Room] == nil {
            h.Rooms[c.Room] = make(map[string]*Client)
        }
        h.Rooms[c.Room][c.UserID] = c
    }

    log.Printf("✅ %s connected to [%s]", c.Username, c.Room)

    safeSend(c, mustJSON(ChatMessage{
        Type:      "system",
        Content:   "Welcome to room [" + c.Room + "]!",
        Timestamp: time.Now().Format("15:04:05"),
    }))
}

func safeSend(c *Client, msg []byte) {
	if c == nil || c.send == nil {
		return
	}
	select {
	case c.send <- msg:
	default:
		log.Printf("[WARN] send channel full or closed for %s", c.Username)
	}
}
func (h *ChatHub) unregister(c *Client) {
    h.mu.Lock()

    current, ok := h.Clients[c.UserID]
    if !ok {
        h.mu.Unlock()
        return
    }

    // Chỉ xử lý nếu đây đúng là connection đang được track
    // (tránh trường hợp connection cũ bị unregister sau khi đã bị replace)
    if current != c {
        // Connection này đã bị replace bởi connection mới → chỉ xóa khỏi room
        if r, ok := h.Rooms[c.Room]; ok {
            delete(r, c.UserID)
            if len(r) == 0 {
                delete(h.Rooms, c.Room)
            }
        }
        h.mu.Unlock()
        log.Printf("⚠️  Stale unregister ignored for %s", c.Username)
        return
    }

    isPrivate := strings.HasPrefix(c.Room, "private_")

    if !isPrivate {
        // Xóa khỏi public room
        if r, ok := h.Rooms[c.Room]; ok {
            delete(r, c.UserID)
            if len(r) == 0 {
                delete(h.Rooms, c.Room)
            }
        }
        // KHÔNG xóa khỏi h.Clients và h.UsernameToUserID
        // vì private connection vẫn đang dùng cùng key c.UserID
        // → chỉ set nil để biết không còn public conn
    } else {
        // Private connection disconnect → xóa hẳn
        delete(h.Clients, c.UserID)
        delete(h.UsernameToUserID, c.Username)
    }

    room := c.Room
    username := c.Username
    h.mu.Unlock()

    log.Printf("❌ %s left [%s]", username, room)

    if !isPrivate {
        h.broadcast(ChatMessage{
            Type:      "leave",
            Sender:    username,
            Content:   username + " left the room",
            Room:      room,
            Timestamp: time.Now().Format("15:04:05"),
        })
    }
}



func (h *ChatHub) broadcast(msg ChatMessage) {
	h.mu.RLock()
	roomClients, ok := h.Rooms[msg.Room]
	if !ok {
		h.mu.RUnlock()
		return
	}

	for _, client := range roomClients {
		select {
		case client.send <- mustJSON(msg):
		default:
			// client chậm
		}
	}
	h.mu.RUnlock()
}

func (h *ChatHub) sendDM(sender *Client, msg ChatMessage) {
	h.mu.RLock()

	userID, exists := h.UsernameToUserID[msg.Recipient]
	if !exists {
		h.mu.RUnlock()
		sender.send <- mustJSON(ChatMessage{
			Type:    "dm_ack",
			ID:      msg.ID,
			Content: "failed",
		})
		return
	}

	receiver, ok := h.Clients[userID]
	h.mu.RUnlock()

	if !ok || receiver == nil {
		sender.send <- mustJSON(ChatMessage{
			Type:    "dm_ack",
			ID:      msg.ID,
			Content: "failed",
		})
		return
	}

	msg.Type = "dm"

	select {
	case receiver.send <- mustJSON(msg):
	default:
	}

	sender.send <- mustJSON(ChatMessage{
		Type:    "dm_ack",
		ID:      msg.ID,
		Content: "delivered",
	})
}

func (h *ChatHub) switchRoom(c *Client, newRoom string) {
	if newRoom == "" || newRoom == c.Room {
		return
	}

	oldRoom := c.Room

	h.mu.Lock()
	delete(h.Rooms[oldRoom], c.UserID)
	if len(h.Rooms[oldRoom]) == 0 {
		delete(h.Rooms, oldRoom)
	}

	if h.Rooms[newRoom] == nil {
		h.Rooms[newRoom] = make(map[string]*Client)
	}
	h.Rooms[newRoom][c.UserID] = c
	c.Room = newRoom
	h.mu.Unlock()

	h.broadcast(ChatMessage{
		Type:      "join",
		Sender:    c.Username,
		Content:   c.Username + " joined the room",
		Room:      newRoom,
		Timestamp: time.Now().Format("15:04:05"),
	})
}

func (h *ChatHub) sendRoomList(c *Client) {
	h.mu.RLock()
	rooms := make([]string, 0, len(h.Rooms))
	for r := range h.Rooms {
		rooms = append(rooms, r)
	}
	h.mu.RUnlock()

	c.send <- mustJSON(ChatMessage{
		Type:      "list_rooms",
		Content:   strings.Join(rooms, ", "),
		Timestamp: time.Now().Format("15:04:05"),
	})
}

func (h *ChatHub) sendUserList(c *Client) {
	h.mu.RLock()
	users := make([]string, 0)
	if room, ok := h.Rooms[c.Room]; ok {
		for _, u := range room {
			users = append(users, u.Username)
		}
	}
	h.mu.RUnlock()

	c.send <- mustJSON(ChatMessage{
		Type:      "user_list",
		Content:   strings.Join(users, ", "),
		Timestamp: time.Now().Format("15:04:05"),
	})
}

// ================= UTILS =================

func normalizeRoom(r string) string {
	r = strings.TrimSpace(strings.ToLower(r))
	if r == "" {
		return "general"
	}
	return r
}

func mustJSON(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

// ================= HANDLER =================

func (h *ChatHub) HandleWS(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	userID := r.URL.Query().Get("userid")
	room := normalizeRoom(r.URL.Query().Get("room"))

	if username == "" || userID == "" {
		http.Error(w, "username and userid are required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	client := &Client{
		Conn:     conn,
		UserID:   userID,
		Username: username,
		Room:     room,
		send:     make(chan []byte, 256),
		hub:      h,
	}

	h.Register <- client

	go client.writePump()
	go client.readPump()
}
