package tcp

import (
	"encoding/json"
	"log"
	"net"
	"sync"
)


type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type ProgressUpdate struct {
	UserID         string `json:"user_id"`
	MangaID        string `json:"manga_id"`
	CurrentChapter int    `json:"current_chapter"`
	Status         string `json:"status"`
	Timestamp      int64  `json:"timestamp"`
}

type AuthMessage struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

type Client struct {
	Conn     net.Conn
	UserID   string
	Username string
}

type ProgressSyncServer struct {
	Port        string
	Clients     map[string]*Client  
	Broadcast   chan ProgressUpdate
	Register    chan *Client
	Unregister  chan *Client
	mu          sync.Mutex
}

func NewProgressSyncServer(port string) *ProgressSyncServer {
	return &ProgressSyncServer{
		Port:       port,
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan ProgressUpdate, 100),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}




func (s *ProgressSyncServer) HandleMessages() {
	for {
		select {
		case client := <-s.Register:
			s.mu.Lock()
			s.Clients[client.UserID] = client
			s.mu.Unlock()
			log.Printf("✓ Can not connect client: %s (%s)\n", client.Username, client.UserID)

		case client := <-s.Unregister:
			s.mu.Lock()
			delete(s.Clients, client.UserID)
			s.mu.Unlock()
			client.Conn.Close()
			log.Printf("✗ Client disconnected: %s\n", client.Username)

		case update := <-s.Broadcast:
			s.mu.Lock()
			for _, client := range s.Clients {
				msg := Message{
					Type:    "progress_update",
					Payload: update,
				}
				data, _ := json.Marshal(msg)
				data = append(data, '\n')

				_, err := client.Conn.Write(data)
				if err != nil {
					log.Println("Error send to client:", client.Username)
				}
			}
			s.mu.Unlock()
		}
	}
}


func (s *ProgressSyncServer) HandleClient(conn net.Conn) {
	client := &Client{Conn: conn}
	defer func() {
		if client.UserID != "" {
			s.Unregister <- client
		} else {
			conn.Close()
		}
	}()

	welcome := Message{Type: "welcome", Payload: "TCP connect success!"}
	data, _ := json.Marshal(welcome)
	conn.Write(append(data, '\n'))

	decoder := json.NewDecoder(conn)
	for {
		var msg Message
		err := decoder.Decode(&msg)
		if err != nil {
			break
		}

		switch msg.Type {
		case "auth":
			authData, _ := json.Marshal(msg.Payload)
			var auth AuthMessage
			json.Unmarshal(authData, &auth)

			client.UserID = auth.UserID
			client.Username = auth.Username
			s.Register <- client

			resp := Message{Type: "auth_success", Payload: "Success auth!"}
			data, _ := json.Marshal(resp)
			conn.Write(append(data, '\n'))

		case "progress_update":
			updateData, _ := json.Marshal(msg.Payload)
			var update ProgressUpdate
			json.Unmarshal(updateData, &update)
			update.UserID = client.UserID

			s.Broadcast <- update
			log.Printf("📡 Broadcasting: %s đọc %s chapter %d\n",
				client.Username, update.MangaID, update.CurrentChapter)
		}
	}
}

func (s *ProgressSyncServer) BroadcastProgress(update ProgressUpdate) {
	s.Broadcast <- update
}