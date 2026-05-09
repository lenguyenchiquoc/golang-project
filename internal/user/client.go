// client.go
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	grpcserver "managahub/internal/grpc"
	// pb "managahub/pkg/proto/managahub/pkg/proto"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/golang-jwt/jwt/v4"
)

const (
	HTTP_SERVER = "http://localhost:8080"
	TCP_SERVER  = "localhost:9090"
	UDP_SERVER  = "localhost:9091"
	WS_SERVER   = "ws://localhost:8080/ws"
)

// ====================== STRUCTS ======================

type TCPMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type ProgressUpdate struct {
	UserID         string `json:"user_id"`
	MangaID        string `json:"manga_id"`
	CurrentChapter int    `json:"current_chapter"`
	Status         string `json:"status"`
	Timestamp      int64  `json:"timestamp"`
}

type UDPNotification struct {
	Type       string `json:"type"`
	MangaID    string `json:"manga_id"`
	MangaTitle string `json:"manga_title"`
	Message    string `json:"message"`
	Timestamp  int64  `json:"timestamp"`
}

type WSMessage struct {
	ID        string `json:"id,omitempty"`
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Room      string `json:"room,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type Session struct {
	Token         string
	UserID        string
	Username      string
	TCPConn       net.Conn
	UDPConn       *net.UDPConn
	WSPrivateConn *websocket.Conn
	mu            sync.Mutex
}

type ChatClient struct {
	conn      *websocket.Conn
	send      chan WSMessage
	done      chan struct{}
	username  string
	room      string
	closeOnce sync.Once
}

type LocalMessage struct {
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

var session = &Session{}
var availableRooms = []string{"general", "manga", "gaming"}
var pendingDM = make(map[string]string)
var pendingMu sync.Mutex
var userListChan = make(chan []string)
var mangaGRPC *grpcserver.MangaGRPCClient

// ====================== MAIN ======================

func main() {
	fmt.Println("╔════════════════════════════════════╗")
	fmt.Println("║        MANGAHUB CLIENT             ║")
	fmt.Println("╚════════════════════════════════════╝")
	getToken := func() string {
		return session.Token
	}

	onAuthFailed := func() {
		doLogout()
	}

	mangaGRPC = grpcserver.NewMangaGRPCClient("localhost:9092", getToken, onAuthFailed)
	defer mangaGRPC.Close()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		if session.Token == "" {
			showAuthMenu(scanner)
		} else {
			showMainMenu(scanner)
		}
	}
}
func startTokenWatcher(token string) { 
    go func() {
        ticker := time.NewTicker(15 * time.Second)
        defer ticker.Stop()

        for {
            <-ticker.C
            if session.Token == "" || session.Token != token {
                return
            }

            if isTokenExpired() {
                fmt.Println("\n❌ Session token is expired!")
                doLogout()
                fmt.Print("Press enter to menu")
                return
            }
        }
    }()
}

// ====================== AUTH MENU ======================
func checkAuthError(resp map[string]interface{}) bool {
	errCode, ok := resp["error"].(string)
	if !ok {
		return false
	}

	switch errCode {
	case "AUTH_TOKEN_EXPIRED":
		fmt.Println("❌ Session expired! Please login again.")
		doLogout()
		return true
	case "AUTH_MISSING_HEADER", "AUTH_INVALID_FORMAT", "AUTH_TOKEN_INVALID":
		fmt.Println("❌ Authentication error! Please login again.")
		doLogout()
		return true
	}
	return false
}

func showAuthMenu(scanner *bufio.Scanner) {
	fmt.Println("\n=== LOGIN / REGISTER ===")
	fmt.Println("1. Login")
	fmt.Println("2. Register")
	fmt.Println("0. Exit")
	fmt.Print("Choice: ")
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())

	switch choice {
	case "1":
		doLogin(scanner)
	case "2":
		doRegister(scanner)
	case "0":
		fmt.Println("Goodbye!")
		os.Exit(0)
	default:
		fmt.Println("❌ Invalid choice")
	}
}

func showMainMenu(scanner *bufio.Scanner) {

	fmt.Printf("\n=== WELCOME %s ===\n", strings.ToUpper(session.Username))
	fmt.Println("1. Search manga")
	fmt.Println("2. View my library")
	fmt.Println("3. Add manga to library")
	fmt.Println("4. Chat")
	fmt.Println("0. Logout")
	fmt.Print("Choice: ")
	if session.Token == "" {
        return 
    }
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())
	
	switch choice {
	case "1":
		doSearchManga(scanner)
	case "2":
		doGetLibrary(scanner)
	case "3":
		doAddToLibrary(scanner)
	case "4":
		chatMenu(scanner)
	case "0":
		doLogout()
	default:
		fmt.Println("❌ Invalid choice")
	}
}


func chatMenu(scanner *bufio.Scanner) {
	for {
		unreadMu.Lock()
		count := unreadDM
		unreadMu.Unlock()
		fmt.Println("\n=== CHAT MENU ===")
		if count > 0 {
			fmt.Printf("1. View DM history 🔴 %d unread\n", count)
		} else {
			fmt.Println("1. View DM history")
		}
		fmt.Println("2. Join room")
		fmt.Println("3. Send DM")
		fmt.Println("0. Back")
		fmt.Print("Choice: ")
		if session.Token == "" {
			return 
		}
		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			unreadMu.Lock()
			unreadDM = 0
			unreadMu.Unlock()
			showDMHistory()
		case "2":
			joinRoomFlow(scanner)
		case "3":
			sendDMFlow(scanner)
		case "0":
			return
		default:
			fmt.Println("❌ Invalid choice")
		}
	}
}


func joinRoomFlow(scanner *bufio.Scanner) {
	fmt.Println("\nAvailable rooms:")
	for i, r := range availableRooms {
		fmt.Printf("%d. %s\n", i+1, r)
	}
	fmt.Print("Enter room name or number: ")
	if session.Token == "" {
        return 
    }	
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	if input == "" {
		fmt.Println("❌ Room name cannot be empty")
		return
	}

	room := input
	for i, r := range availableRooms {
		if input == fmt.Sprintf("%d", i+1) {
			room = r
			break
		}
	}

	urlStr := fmt.Sprintf("%s?username=%s&userid=%s&room=%s",
		WS_SERVER, session.Username, session.UserID, room)

	conn, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
	if err != nil {
		fmt.Printf("❌ Cannot connect to WebSocket: %v\n", err)
		return
	}

	client := newChatClient(conn, session.Username, room)

	fmt.Printf("\n=== Joined room: %s ===\n", room)
	fmt.Println("Commands: /users | /rooms | /switch <room> | /leave | /exit | /help")

	go client.writeLoop()
	go client.readLoop()
	roomLoop(scanner, client)
}

func roomLoop(scanner *bufio.Scanner, client *ChatClient) {
	for {
		if session.Token == "" {
			client.close()
			return
		}
		fmt.Printf("[%s] > ", client.room)
		if !scanner.Scan() {
			client.close()
			return
		}

		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}

		if session.Token == "" {
			client.close()
			return
		}

		switch {
		case text == "/help":
			fmt.Println("  /dm <username> <message>   → Send private message")
			fmt.Println("  /users                     → List users in current room")
			fmt.Println("  /rooms                     → List all active rooms")
			fmt.Println("  /switch <room>             → Switch to another room")
			fmt.Println("  /leave                     → Leave current room")
			fmt.Println("  /exit                      → Exit chat completely")

		case text == "/leave":
			client.close()
			fmt.Println("👋 Left the room")
			return

		case text == "/exit":
			client.close()
			fmt.Println("👋 Exited chat")
			os.Exit(0)

		case text == "/users":
			client.send <- WSMessage{Type: "list_users"}

		case text == "/rooms":
			client.send <- WSMessage{Type: "list_rooms"}

		case strings.HasPrefix(text, "/switch "):
			newRoom := strings.TrimSpace(text[8:])
			if newRoom == "" {
				fmt.Println("Usage: /switch <room>")
				continue
			}
			client.room = newRoom
			client.send <- WSMessage{
				Type:    "join_room",
				Content: newRoom,
			}
			fmt.Printf("=== Switching to room: %s ===\n", newRoom)

		default:
			// Normal chat message
			client.send <- WSMessage{
				Type:    "chat",
				Content: text,
			}
		}
	}
}

// ====================== SEND DM DIRECTLY ======================
var lastUserList []string

func sendDMFlow(scanner *bufio.Scanner) {
	fmt.Println("1. Choose online user")
	fmt.Println("2. Enter username manually")
	fmt.Print("Choice: ")

	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())

	var recipient string
	if session.Token == "" {
            return 
    }
	switch choice {

	case "1":

		if session.WSPrivateConn == nil {
			doListenWS()
			time.Sleep(700 * time.Millisecond) 
		}
		if session.WSPrivateConn == nil {
			fmt.Println("❌ Not connected to chat server")
			return
		}

		session.WSPrivateConn.WriteMessage(
			websocket.TextMessage,
			[]byte(`{"type":"list_all_users"}`),
		)

		var users []string

		select {
		case users = <-userListChan:
		case <-time.After(2 * time.Second):
			fmt.Println("❌ Timeout getting user list")
			return
		}

		if len(users) == 0 {
			fmt.Println("❌ No users online")
			return
		}

		fmt.Println("\n=== ONLINE USERS ===")
		for i, u := range users {
			fmt.Printf("%d. %s\n", i+1, u)
		}

		fmt.Print("Select user: ")
		scanner.Scan()

		var idx int
		fmt.Sscanf(scanner.Text(), "%d", &idx)

		if idx < 1 || idx > len(users) {
			fmt.Println("❌ Invalid choice")
			return
		}

		recipient = users[idx-1]

	case "2":
		fmt.Print("To (username): ")
		scanner.Scan()
		recipient = strings.TrimSpace(scanner.Text())

		if recipient == "" {
			fmt.Println("❌ Username cannot be empty")
			return
		}

	default:
		fmt.Println("❌ Invalid choice")
		return
	}

	// ================= SEND =================

	fmt.Print("Message: ")
	scanner.Scan()
	msgContent := strings.TrimSpace(scanner.Text())
	if msgContent == "" {
		fmt.Println("❌ Message cannot be empty")
		return
	}

	msgID := uuid.New().String()

	pendingMu.Lock()
	pendingDM[msgID] = "sending"
	pendingMu.Unlock()

	data, _ := json.Marshal(WSMessage{
		ID:        msgID,
		Type:      "dm",
		Recipient: recipient,
		Content:   msgContent,
	})

	err := session.WSPrivateConn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		fmt.Println("❌ Failed to send DM")
		return
	}

	fmt.Printf("📨 Sending DM to %s...\n", recipient)
}

// ====================== CHAT CLIENT ======================

func newChatClient(conn *websocket.Conn, username, room string) *ChatClient {
	c := &ChatClient{
		conn:     conn,
		send:     make(chan WSMessage, 32),
		done:     make(chan struct{}),
		username: username,
		room:     room,
	}

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			case <-c.done:
				return
			}
		}
	}()

	return c
}

var unreadDM int
var unreadMu sync.Mutex

func (c *ChatClient) readLoop() {
	defer c.close()

	for {

		if session.Token == "" {
			return
		}
		c.conn.SetPongHandler(func(string) error {
			c.conn.SetReadDeadline(time.Now().Add(70 * time.Second))
			return nil
		})
		_, rawMsg, err := c.conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				fmt.Println("\n❌ Disconnected from chat server")
			}
			return
		}

		var msg WSMessage
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			continue
		}

		if session.Token == "" {
			return
		}

		switch msg.Type {
		case "chat":
			fmt.Printf("\n💬 [%s] %s: %s\n", msg.Timestamp, msg.Sender, msg.Content)
		case "dm":
			fmt.Printf("\n📩 [DM] %s: %s\n", msg.Sender, msg.Content)
			unreadMu.Lock()
			unreadDM++
			unreadMu.Unlock()
		case "system":
			fmt.Printf("\n● %s\n", msg.Content)
		case "join":
			fmt.Printf("\n👋 %s joined the room\n", msg.Sender)
		case "leave":
			fmt.Printf("\n👋 %s left the room\n", msg.Sender)
		case "list_rooms":
			fmt.Printf("\n📦 Active rooms: %s\n", msg.Content)
		case "user_list":
			fmt.Printf("\n👥 Users in room: %s\n", msg.Content)
		case "error":
			fmt.Printf("\n❌ %s\n", msg.Content)
		}

		fmt.Printf("[%s] > ", c.room)
	}
}

func (c *ChatClient) writeLoop() {
	defer c.close()

	for {
		select {
		case msg := <-c.send:
			data, err := json.Marshal(msg)
			if err != nil {
				continue
			}
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				fmt.Println("\n❌ Failed to send message")
				return
			}
		case <-c.done:
			return
		}
	}
}

func (c *ChatClient) close() {
	c.closeOnce.Do(func() {
		close(c.done)
		if c.conn != nil {
			c.conn.Close()
		}
	})
}

// ====================== DM HISTORY ======================

func showDMHistory() {
	history := loadDMHistory()
	if len(history) == 0 {
		fmt.Println("📭 No DM history found")
		return
	}

	fmt.Println("\n=== DM HISTORY ===")
	for _, m := range history {
		if m.Sender == session.Username {
			fmt.Printf("➡️  To %s [%s]: %s\n", m.Recipient, m.Timestamp, m.Content)
		} else {
			fmt.Printf("⬅️  From %s [%s]: %s\n", m.Sender, m.Timestamp, m.Content)
		}
	}
}

// ====================== AUTH FUNCTIONS ======================

func doRegister(scanner *bufio.Scanner) {
	fmt.Print("Username: ")
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())

	fmt.Print("Email: ")
	scanner.Scan()
	email := strings.TrimSpace(scanner.Text())

	fmt.Print("Password: ")
	scanner.Scan()
	password := strings.TrimSpace(scanner.Text())

	fmt.Print("Re-enter password: ")
	scanner.Scan()
	rePassword := strings.TrimSpace(scanner.Text())

	body := map[string]string{
		"username":   username,
		"email":      email,
		"password":   password,
		"RePassword": rePassword,
	}

	resp, err := httpPost("/auth/register", body, "")
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	if msg, ok := resp["message"]; ok {
		fmt.Println("✅", msg)
	} else if errMsg, ok := resp["error"]; ok {
		fmt.Println("❌", errMsg)
	}
}

func doLogin(scanner *bufio.Scanner) {
	fmt.Print("Username: ")
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())

	fmt.Print("Password: ")
	scanner.Scan()
	password := strings.TrimSpace(scanner.Text())

	body := map[string]string{
		"username": username,
		"password": password,
	}

	resp, err := httpPost("/auth/login", body, "")
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	if token, ok := resp["token"]; ok {
		session.Token = token.(string)
		session.UserID = resp["user_id"].(string)
		session.Username = resp["username"].(string)
		fmt.Printf("✅ Login successful! Welcome %s\n", session.Username)
		go doConnectTCP()
		go doRegisterUDP()
		go doListenWS()
		startTokenWatcher(session.Token)
	} else if errMsg, ok := resp["error"]; ok {
		fmt.Println("❌", errMsg)
	}
}

func doLogout() {
	session.mu.Lock()

	_, err := httpPost("/auth/logout", nil, session.Token)
	if err != nil {
		fmt.Println("❌ Logout API failed:", err)
	}

	if session.TCPConn != nil {
		session.TCPConn.Close()
	}
	if session.UDPConn != nil {
		doUnregisterUDP()
		session.UDPConn.Close()
	}
	if session.WSPrivateConn != nil {
		session.WSPrivateConn.Close()
	}
	session.mu.Unlock()

	fmt.Println("✅ Logged out successfully!")
	session = &Session{}
}



// ====================== MANGA FUNCTIONS ======================

func doSearchManga(scanner *bufio.Scanner) {
	if session.Token == "" {
            return 
    }
	fmt.Print("Results per page (default 10): ")
	scanner.Scan()
	limitStr := strings.TrimSpace(scanner.Text())

	limit := int32(10)
	if limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}

	// query
	fmt.Print("Enter manga name or author: ")
	scanner.Scan()
	query := strings.TrimSpace(scanner.Text())

	// genre (multi)
	fmt.Print("Enter genres (comma separated, optional): ")
	scanner.Scan()
	genreInput := strings.TrimSpace(scanner.Text())

	var genres []string
	if genreInput != "" {
		for _, g := range strings.Split(genreInput, ",") {
			genres = append(genres, strings.TrimSpace(g))
		}
	}

	// status
	fmt.Print("Enter status (ongoing/completed, optional): ")
	scanner.Scan()
	status := strings.TrimSpace(scanner.Text())

	resp, err := mangaGRPC.SearchManga(query, genres, status, 1, limit)
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	if len(resp.Mangas) == 0 {
		fmt.Println("No results found!")
		return
	}

	fmt.Printf("\n📚 Found %d results:\n", resp.Total)
	fmt.Println("─────────────────────────────────────────────────────")

	for i, m := range resp.Mangas {
		fmt.Printf("%d. [%s] %s - %s (%d chapters) [%s]\n",
			i+1, m.Id, m.Title, m.Author, m.TotalChapters, m.Status,
		)

		// in genres
		if len(m.Genres) > 0 {
			fmt.Print("   Genres: ")
			for j, g := range m.Genres {
				fmt.Print(g.Name)
				if j < len(m.Genres)-1 {
					fmt.Print(", ")
				}
			}
			fmt.Println()
		}
	}

	fmt.Println("─────────────────────────────────────────────────────")

	fmt.Print("Choose manga (number, 0 to cancel): ")
	scanner.Scan()
	var choice int
	fmt.Sscanf(scanner.Text(), "%d", &choice)

	if choice <= 0 || choice > len(resp.Mangas) {
		fmt.Println("❌ Cancel")
		return
	}

	selected := resp.Mangas[choice-1]

	viewMangaDetailAndRate(scanner, selected.Id)
}

func viewMangaDetailAndRate(scanner *bufio.Scanner, mangaID string) {
	if session.Token == "" {
            return 
    }	
	m, err := mangaGRPC.GetManga(mangaID)
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	fmt.Println("\n📖 Manga Detail")
	fmt.Println("────────────────────────────────────────")
	fmt.Println("Title:", m.Title)
	fmt.Println("Author:", m.Author)
	fmt.Println("Status:", m.Status)
	fmt.Println("Chapters:", m.TotalChapters)

	fmt.Print("Genres: ")
	for i, g := range m.Genres {
		fmt.Print(g.Name)
		if i < len(m.Genres)-1 {
			fmt.Print(", ")
		}
	}
	fmt.Println()

	fmt.Printf("⭐ Rating: %.1f (%d votes)\n", m.AverageRating, m.RatingCount)
	fmt.Println("────────────────────────────────────────")

	fmt.Print("Do you want to rate this manga? (y/n): ")
	scanner.Scan()
	ans := strings.TrimSpace(scanner.Text())

	if ans != "y" && ans != "Y" {
		return
	}

	fmt.Print("Enter rating (1-10): ")
	scanner.Scan()
	var rating int32
	fmt.Sscanf(scanner.Text(), "%d", &rating)

	resp, err := mangaGRPC.RateManga(session.UserID, mangaID, rating)
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	if !resp.Success {
		fmt.Println("❌", resp.Message)
		return
	}

	fmt.Printf("✅ Rated successfully! New average: %.1f (%d votes)\n",
		resp.Average, resp.Count)
}

func doGetLibrary(scanner *bufio.Scanner) {
	if session.Token == "" {
            return 
    }
	resp, err := httpGet("/users/library", session.Token)
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}
	if checkAuthError(resp) {
		return
	}

	items, ok := resp["library"].([]interface{})
	if !ok || len(items) == 0 {
		fmt.Println("📚 Library is empty!")
		return
	}

	fmt.Printf("\n📚 %s's library (%v manga):\n", session.Username, resp["total"])
	fmt.Println("──────────────────────────────────────────────────────────")
	var lst []string 
	for i, item := range items {
		m := item.(map[string]interface{})
		fmt.Printf("%d. %s\n", i+1, m["title"])
		fmt.Printf("   ID       : %s\n", m["manga_id"])
		lst = append(lst, m["manga_id"].(string))
		fmt.Printf("   Author   : %s\n", m["author"])
		fmt.Printf("   Progress : Chapter %v / %v\n", m["current_chapter"], m["total_chapters"])
		fmt.Printf("   Status   : %s\n", m["status"])
		fmt.Println("──────────────────────────────────────────────────────────")
	}

	fmt.Println("─────────────────────────────────────────────────────")
	fmt.Print("Choose manga to update progress (number, 0 to cancel): ")
	scanner.Scan()
	var choice int
	fmt.Sscanf(scanner.Text(), "%d", &choice)

	if choice <= 0 || choice > len(lst) {
		fmt.Println("❌ Cancel b/c invalid input")
		return
	}
	choose := lst[choice-1]
	
	doUpdateProgress(scanner, choose)

}

func doAddToLibrary(scanner *bufio.Scanner) {
	if session.Token == "" {
            return 
    }
	fmt.Print("Enter Manga ID: ")
	scanner.Scan()
	mangaID := strings.TrimSpace(scanner.Text())

	fmt.Println("Status:")
	fmt.Println("1. reading")
	fmt.Println("2. completed")
	fmt.Println("3. plan_to_read")
	fmt.Print("Choice: ")
	scanner.Scan()
	statusChoice := strings.TrimSpace(scanner.Text())

	statusMap := map[string]string{"1": "reading", "2": "completed", "3": "plan_to_read"}
	status := statusMap[statusChoice]
	if status == "" {
		fmt.Println("❌ Invalid choice")
		return
	}

	body := map[string]interface{}{
		"manga_id":        mangaID,
		"status":          status,
		"current_chapter": 0,
	}

	resp, err := httpPost("/users/library", body, session.Token)
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	if checkAuthError(resp) {
		return
	}

	if msg, ok := resp["message"]; ok {
		fmt.Println("✅", msg)
	} else {
		fmt.Println("❌", resp["error"])
	}
}

func doUpdateProgress(scanner *bufio.Scanner, choose string) {
	if session.Token == "" {
            return 
    }
	mangaID := strings.TrimSpace(choose)

	fmt.Print("Current chapter: ")
	scanner.Scan()
	chapterStr := strings.TrimSpace(scanner.Text())

	var chapter int
	fmt.Sscanf(chapterStr, "%d", &chapter)

	fmt.Println("Status (Enter to keep current):")
	fmt.Println("1. reading")
	fmt.Println("2. completed")
	fmt.Println("3. plan_to_read")
	fmt.Print("Choice: ")
	scanner.Scan()
	statusChoice := strings.TrimSpace(scanner.Text())

	statusMap := map[string]string{"1": "reading", "2": "completed", "3": "plan_to_read"}
	status := statusMap[statusChoice]

	body := map[string]interface{}{
		"manga_id":        mangaID,
		"current_chapter": chapter,
		"status":          status,
	}

	resp, err := httpPut("/users/progress", body, session.Token)
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	if checkAuthError(resp) {
		return
	}

	if msg, ok := resp["message"]; ok {
		fmt.Printf("✅ %s - Chapter %d\n", msg, chapter)
	} else {
		fmt.Println("❌", resp["error"])
	}
}

// ====================== NETWORK FUNCTIONS ======================

func doConnectTCP() {
	if session.TCPConn != nil {
		return
	}

	conn, err := net.Dial("tcp", TCP_SERVER)
	if err != nil {
		fmt.Println("❌ Cannot connect to TCP server:", err)
		return
	}
	session.TCPConn = conn

	auth := map[string]interface{}{
		"type": "auth",
		"payload": map[string]string{
			"user_id":  session.UserID,
			"username": session.Username,
		},
	}
	data, _ := json.Marshal(auth)
	conn.Write(append(data, '\n'))

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			var msg TCPMessage
			if err := json.Unmarshal([]byte(scanner.Text()), &msg); err != nil {
				continue
			}
			switch msg.Type {
			case "progress_update":
				var update ProgressUpdate
				json.Unmarshal(msg.Payload, &update)
				t := time.Unix(update.Timestamp, 0).Format("15:04:05")
				fmt.Printf("\n📡 [TCP][%s] %s updated %s to chapter %d\n",
					t, update.UserID, update.MangaID, update.CurrentChapter)
				fmt.Print("Choice: ")
			case "auth_success":
				fmt.Println("✅ TCP connected!")
			}
		}
		fmt.Println("❌ TCP connection lost!")
		session.TCPConn = nil
	}()
}

func doRegisterUDP() {
	serverAddr, err := net.ResolveUDPAddr("udp", UDP_SERVER)
	if err != nil {
		fmt.Println("❌ Cannot resolve UDP address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("❌ Cannot connect to UDP server:", err)
		return
	}
	session.UDPConn = conn

	packet := map[string]string{
		"type":    "register",
		"user_id": session.UserID,
	}
	data, _ := json.Marshal(packet)
	conn.Write(data)

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				break
			}

			var notification UDPNotification
			if err := json.Unmarshal(buf[:n], &notification); err != nil {
				continue
			}

			switch notification.Type {
			case "register_success":
				fmt.Println("✅ UDP notifications registered!")
			case "chapter_release":
				t := time.Unix(notification.Timestamp, 0).Format("15:04:05")
				fmt.Printf("\n🔔 [UDP][%s] %s - %s\n",
					t, notification.MangaTitle, notification.Message)
				fmt.Print("Choice: ")
			}
		}
		fmt.Println("❌ UDP connection lost!")
		session.UDPConn = nil
	}()
}

func doUnregisterUDP() {
	if session.UDPConn == nil {
		return
	}
	packet := map[string]string{
		"type":    "unregister",
		"user_id": session.UserID,
	}
	data, _ := json.Marshal(packet)
	session.UDPConn.Write(data)
}

// HTTP Helpers
func httpPost(path string, body interface{}, token string) (map[string]interface{}, error) {
	return httpRequest("POST", path, body, token)
}

func httpPut(path string, body interface{}, token string) (map[string]interface{}, error) {
	return httpRequest("PUT", path, body, token)
}

func httpGet(path string, token string) (map[string]interface{}, error) {
	return httpRequest("GET", path, nil, token)
}

func httpRequest(method, path string, body interface{}, token string) (map[string]interface{}, error) {
	var reqBody io.Reader
	if body != nil {
		data, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(data)
	}

	req, err := http.NewRequest(method, HTTP_SERVER+path, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

// ====================== DM STORAGE ======================

func saveDM(msg LocalMessage) {
	file := getConversationFile(msg.Sender, msg.Recipient)

	var history []LocalMessage
	data, err := os.ReadFile(file)
	if err == nil {
		json.Unmarshal(data, &history)
	}

	history = append(history, msg)

	newData, _ := json.MarshalIndent(history, "", "  ")
	os.WriteFile(file, newData, 0644)
}

func loadDMHistory() []LocalMessage {
	var allMessages []LocalMessage
	files, err := os.ReadDir(".")
	if err != nil {
		return allMessages
	}

	for _, f := range files {
		name := f.Name()
		if !strings.HasPrefix(name, "dm_") || !strings.HasSuffix(name, ".json") {
			continue
		}

		data, err := os.ReadFile(name)
		if err != nil {
			continue
		}

		var history []LocalMessage
		if err := json.Unmarshal(data, &history); err != nil {
			continue
		}

		for _, msg := range history {
			if msg.Sender == session.Username || msg.Recipient == session.Username {
				allMessages = append(allMessages, msg)
			}
		}
	}
	return allMessages
}

func getConversationFile(a, b string) string {
	if a > b {
		a, b = b, a
	}
	return fmt.Sprintf("dm_%s_%s.json", a, b)
}

// func doListenWS() {
// 	privateRoom := "private_" + session.UserID
// 	wsURL := fmt.Sprintf("%s?username=%s&userid=%s&room=%s",
// 		WS_SERVER, session.Username, session.UserID, privateRoom)

// 	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
// 	if err != nil {
// 		return
// 	}

// 	session.WSPrivateConn = conn

// 	go func() {
// 		for {
// 			_, raw, err := conn.ReadMessage()
// 			if err != nil {
// 				session.WSPrivateConn = nil
// 				return
// 			}

// 			var msg WSMessage
// 			if json.Unmarshal(raw, &msg) != nil {
// 				continue
// 			}

// 			switch msg.Type {
// 			case "dm":
// 				saveDM(LocalMessage{
// 					Type:      "dm",
// 					Sender:    msg.Sender,
// 					Recipient: session.Username,
// 					Content:   msg.Content,
// 					Timestamp: msg.Timestamp,
// 				})
// 				unreadMu.Lock()
// 				unreadDM++
// 				unreadMu.Unlock()

// 			case "error":
// 				fmt.Printf("\n❌ DM error: %s\n", msg.Content)
// 				fmt.Print("Choice: ")

// 			case "dm_ack":
// 				pendingMu.Lock()
// 				_, exists := pendingDM[msg.ID]
// 				if exists {
// 					delete(pendingDM, msg.ID)
// 				}
// 				pendingMu.Unlock()

// 				if msg.Content == "delivered" {
// 					fmt.Println("\n✅ DM delivered")
// 				} else {
// 					fmt.Println("\n❌ DM failed")
// 				}
// 				fmt.Print("Choice: ")

// 			case "list_all_users":
// 				users := strings.Split(msg.Content, ", ")
// 				userListChan <- users

// 			}
// 		}
// 	}()

// }

func doListenWS() {
	session.mu.Lock()
	if session.WSPrivateConn != nil {
		session.mu.Unlock()
		return
	}
	session.mu.Unlock()

	privateRoom := "private_" + session.UserID
	wsURL := fmt.Sprintf("%s?username=%s&userid=%s&room=%s",
		WS_SERVER, session.Username, session.UserID, privateRoom)

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		fmt.Println("❌ Cannot connect private WS:", err)
		return
	}

	session.mu.Lock()
	session.WSPrivateConn = conn
	session.mu.Unlock()

	fmt.Println("✅ Private WS connected successfully")

	go func() {
		defer func() {
			session.mu.Lock()
			session.WSPrivateConn = nil
			session.mu.Unlock()
		}()

		for {
			if session.Token == "" {
				return
			}

			_, raw, err := conn.ReadMessage()
			if err != nil {
				if session.Token != "" {
					fmt.Println("❌ Private WS disconnected")
				}
				return
			}

			var msg WSMessage
			if err := json.Unmarshal(raw, &msg); err != nil {
				continue
			}

			switch msg.Type {
			case "list_all_users":
				usersStr := strings.TrimSpace(msg.Content)
				users := strings.Split(usersStr, ", ")
				if usersStr == "" {
					users = []string{}
				}
				select {
				case userListChan <- users:
				default:
					
				}

			case "dm":
				saveDM(LocalMessage{
					Type:      "dm",
					Sender:    msg.Sender,
					Recipient: session.Username,
					Content:   msg.Content,
					Timestamp: msg.Timestamp,
				})
				unreadMu.Lock()
				unreadDM++
				unreadMu.Unlock()

			case "error":
				fmt.Printf("\n❌ DM error: %s\n", msg.Content)
				fmt.Print("Choice: ")

			case "dm_ack":
				pendingMu.Lock()
				_, exists := pendingDM[msg.ID]
				if exists {
					delete(pendingDM, msg.ID)
				}
				pendingMu.Unlock()

				if msg.Content == "delivered" {
					fmt.Println("\n✅ DM delivered")
				} else {
					fmt.Println("\n❌ DM failed")
				}
				fmt.Print("Choice: ")
			}
		}
	}()
}

func isTokenExpired() bool {
    if session.Token == "" {
        return true
    }

    // Parse token không verify signature
    token, _, err := new(jwt.Parser).ParseUnverified(session.Token, jwt.MapClaims{})
    if err != nil {
        return true
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return true
    }

    // Check exp
    exp, ok := claims["exp"].(float64)
    if !ok {
        return true
    }

    return time.Now().Unix() > int64(exp)
}
