package udp

import (
	"encoding/json"
	"log"
	"net"
	"sync"
)

type Notification struct {
	Type       string `json:"type"`
	MangaID    string `json:"manga_id"`
	MangaTitle string `json:"manga_title"`
	Message    string `json:"message"`
	Timestamp  int64  `json:"timestamp"`
}

type RegisterPacket struct {
	Type   string `json:"type"`
	UserID string `json:"user_id"`
}

type NotificationServer struct {
	Port    string
	Clients map[string]*net.UDPAddr
	conn    *net.UDPConn
	mu      sync.Mutex
}

func NewNotificationServer(port string) *NotificationServer {
	return &NotificationServer{
		Port:    port,
		Clients: make(map[string]*net.UDPAddr),
	}
}


func (s *NotificationServer) HandlePacket(conn *net.UDPConn, data []byte, clientAddr *net.UDPAddr) {
	var packet RegisterPacket
	err := json.Unmarshal(data, &packet)
	if err != nil {
		log.Println("Error parsing packet:", err)
		return
	}

	switch packet.Type {
	case "register":
		s.mu.Lock()
		s.Clients[packet.UserID] = clientAddr
		s.mu.Unlock()

		log.Printf("Client registered UDP: %s (%s)\n", packet.UserID, clientAddr)

		resp, _ := json.Marshal(map[string]string{
			"type":    "register_success",
			"message": "Successfully registered for notifications!",
		})
		conn.WriteToUDP(resp, clientAddr)

	case "unregister":
		s.mu.Lock()
		delete(s.Clients, packet.UserID)
		s.mu.Unlock()

		log.Printf("Client unregistered UDP: %s\n", packet.UserID)

		resp, _ := json.Marshal(map[string]string{
			"type":    "unregister_success",
			"message": "Successfully unregistered from notifications!",
		})
		conn.WriteToUDP(resp, clientAddr)
	}
}

func (s *NotificationServer) Broadcast(notification Notification) {
	if s.conn == nil {
		log.Println("UDP server not started yet")
		return
	}

	data, err := json.Marshal(notification)
	if err != nil {
		log.Println("Error marshaling notification:", err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	successCount := 0
	for userID, clientAddr := range s.Clients {
		_, err := s.conn.WriteToUDP(data, clientAddr)
		if err != nil {
			log.Printf("Error sending to %s: %v\n", userID, err)
			continue
		}
		successCount++
	}

	log.Printf("UDP Broadcast: %s → %d clients\n", notification.Message, successCount)
}

func (s *NotificationServer) GetClientCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.Clients)
}

func (s *NotificationServer) SetConn(conn *net.UDPConn) {
	s.conn = conn
}