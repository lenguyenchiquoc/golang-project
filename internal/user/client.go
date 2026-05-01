package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	HTTP_SERVER = "http://localhost:8080"
	TCP_SERVER  = "localhost:9090"
	UDP_SERVER  = "localhost:9091"
	WS_SERVER   = "ws://localhost:8080/ws"
)

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
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	Room      string `json:"room,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type Session struct {
	Token    string
	UserID   string
	Username string
	TCPConn  net.Conn
	UDPConn  *net.UDPConn
	WSConn   *websocket.Conn
	WSRoom   string
	mu       sync.Mutex
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

func newChatClient(conn *websocket.Conn, username, room string) *ChatClient {
	return &ChatClient{
		conn:     conn,
		send:     make(chan WSMessage, 10),
		done:     make(chan struct{}),
		username: username,
		room:     room,
	}
}

var session = &Session{}
var availableRooms = []string{"general", "manga", "gaming"}

func main() {
	fmt.Println("╔════════════════════════════════════╗")
	fmt.Println("║        MANGAHUB CLIENT             ║")
	fmt.Println("╚════════════════════════════════════╝")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		if session.Token == "" {
			showAuthMenu(scanner)
		} else {
			showMainMenu(scanner)
		}
	}
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
	fmt.Println("4. Update reading progress")
	fmt.Println("5. Chat")
	fmt.Println("0. Logout")
	fmt.Print("Choice: ")

	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())

	switch choice {
	case "1":
		doSearchManga(scanner)
	case "2":
		doGetLibrary()
	case "3":
		doAddToLibrary(scanner)
	case "4":
		doUpdateProgress(scanner)
	case "5":
		chatMenu(scanner)
	case "0":
		doLogout()
	default:
		fmt.Println("❌ Invalid choice")
	}
}

func chatMenu(scanner *bufio.Scanner) {
	for {
		fmt.Println("\n=== CHAT MENU ===")
		fmt.Println("1. View DM history")
		fmt.Println("2. Join room")
		fmt.Println("3. Send DM")
		fmt.Println("0. Back")
		fmt.Print("Choice: ")

		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			history := loadDMHistory()
			if len(history) == 0 {
				fmt.Println("📭 No DM history")
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
	fmt.Print("Select room (number or name): ")
	scanner.Scan()
	input := strings.TrimSpace(scanner.Text())

	if input == "" {
		fmt.Println("❌ Room name cannot be empty")
		return
	}

	// Chọn theo số hoặc tên
	room := input
	for i, r := range availableRooms {
		if input == fmt.Sprintf("%d", i+1) {
			room = r
			break
		}
	}

	url := fmt.Sprintf("%s?username=%s&room=%s", WS_SERVER, session.Username, room)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Println("❌ Cannot connect:", err)
		return
	}

	client := newChatClient(conn, session.Username, room)
	go client.writeLoop()
	roomLoop(scanner, client)
}

func roomLoop(scanner *bufio.Scanner, client *ChatClient) {
	go client.readLoop()

	fmt.Printf("\n=== ROOM: %s ===\n", client.room)
	fmt.Println("Commands:")
	fmt.Println("  /dm <user> <msg>     → send private message")
	fmt.Println("  /users               → list users in room")
	fmt.Println("  /rooms               → list all rooms")
	fmt.Println("  /switch <room>       → switch room")
	fmt.Println("  /leave               → leave room (back to chat menu)")
	fmt.Println("  /exit                → exit chat completely")
	fmt.Println("  /help                → show commands")
	fmt.Println("──────────────────────────────────────")

	for {
		fmt.Printf("[%s] > ", client.room)
		if !scanner.Scan() {
			return
		}
		text := strings.TrimSpace(scanner.Text())

		if text == "" {
			continue
		}

		switch {
		case text == "/help":
			fmt.Println("  /dm <user> <msg>     → send private message")
			fmt.Println("  /users               → list users in room")
			fmt.Println("  /rooms               → list all rooms")
			fmt.Println("  /switch <room>       → switch room")
			fmt.Println("  /leave               → leave room (back to chat menu)")
			fmt.Println("  /exit                → exit chat completely")

		case text == "/leave":
			client.close()
			fmt.Println("👋 Left room")
			return

		case text == "/exit":
			client.close()
			fmt.Println("👋 Exit chat")
			return

		case strings.HasPrefix(text, "/dm "):
			parts := strings.SplitN(text[4:], " ", 2)
			if len(parts) < 2 {
				fmt.Println("Usage: /dm <user> <msg>")
				continue
			}
			client.send <- WSMessage{
				Type:      "dm",
				Recipient: parts[0],
				Content:   parts[1],
			}
			saveDM(LocalMessage{
				Type:      "dm",
				Sender:    client.username,
				Recipient: parts[0],
				Content:   parts[1],
				Timestamp: time.Now().Format("15:04:05"),
			})
		case text == "/users":
			client.send <- WSMessage{Type: "list_users"}

		case text == "/rooms":
			client.send <- WSMessage{Type: "list_rooms"}

		case strings.HasPrefix(text, "/switch "):
			room := strings.TrimSpace(text[8:])
			if room == "" {
				fmt.Println("Usage: /switch <room>")
				continue
			}
			client.room = room
			client.send <- WSMessage{Type: "join_room", Content: room}
			fmt.Printf("=== ROOM: %s ===\n", room)

		default:
			client.send <- WSMessage{
				Type:    "chat",
				Content: text,
			}
		}
	}
}

func sendDMFlow(scanner *bufio.Scanner) {
	fmt.Print("To (username): ")
	scanner.Scan()
	user := strings.TrimSpace(scanner.Text())
	if user == "" {
		fmt.Println("❌ Username cannot be empty")
		return
	}

	fmt.Print("Message: ")
	scanner.Scan()
	msg := strings.TrimSpace(scanner.Text())
	if msg == "" {
		fmt.Println("❌ Message cannot be empty")
		return
	}

	url := fmt.Sprintf("%s?username=%s&room=general", WS_SERVER, session.Username)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Println("❌ Cannot connect to chat server")
		return
	}
	defer conn.Close()

	time.Sleep(300 * time.Millisecond)

	data, _ := json.Marshal(WSMessage{
		Type:      "dm",
		Recipient: user,
		Content:   msg,
	})
	conn.WriteMessage(websocket.TextMessage, data)

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, raw, err := conn.ReadMessage()
	if err == nil {
		var resp WSMessage
		if json.Unmarshal(raw, &resp) == nil && resp.Type == "error" {
			fmt.Println("❌", resp.Content)
			return
		}
	}

	fmt.Printf("✅ DM sent to %s\n", user)

	saveDM(LocalMessage{
		Type:      "dm",
		Sender:    session.Username,
		Recipient: user,
		Content:   msg,
		Timestamp: time.Now().Format("15:04:05"),
	})
}

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
	} else {
		fmt.Println("❌", resp["error"])
	}
}

func doLogout() {
	if session.TCPConn != nil {
		session.TCPConn.Close()
	}
	if session.UDPConn != nil {
		doUnregisterUDP()
		session.UDPConn.Close()
	}
	if session.WSConn != nil {
		session.WSConn.Close()
	}
	session = &Session{}
	fmt.Println("✅ Logged out!")
}

func doSearchManga(scanner *bufio.Scanner) {
	fmt.Print("Enter manga name or author: ")
	scanner.Scan()
	query := strings.TrimSpace(scanner.Text())

	resp, err := httpGet("/manga?query="+query, session.Token)
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	results, ok := resp["results"].([]interface{})
	if !ok || len(results) == 0 {
		fmt.Println("No results found!")
		return
	}

	fmt.Printf("\n📚 Found %v results:\n", resp["total"])
	fmt.Println("─────────────────────────────────────────────────────")
	for i, r := range results {
		m := r.(map[string]interface{})
		fmt.Printf("%d. [%s] %s - %s (%v chapters) [%s]\n",
			i+1, m["id"], m["title"], m["author"], m["total_chapters"], m["status"],
		)
	}
	fmt.Println("─────────────────────────────────────────────────────")
}

func doGetLibrary() {
	resp, err := httpGet("/users/library", session.Token)
	if err != nil {
		fmt.Println("❌ Error:", err)
		return
	}

	items, ok := resp["library"].([]interface{})
	if !ok || len(items) == 0 {
		fmt.Println("📚 Library is empty!")
		return
	}

	fmt.Printf("\n📚 %s's library (%v manga):\n", session.Username, resp["total"])
	fmt.Println("──────────────────────────────────────────────────────────")
	for i, item := range items {
		m := item.(map[string]interface{})
		fmt.Printf("%d. %s\n", i+1, m["title"])
		fmt.Printf("   ID       : %s\n", m["manga_id"])
		fmt.Printf("   Author   : %s\n", m["author"])
		fmt.Printf("   Progress : Chapter %v / %v\n", m["current_chapter"], m["total_chapters"])
		fmt.Printf("   Status   : %s\n", m["status"])
		fmt.Println("──────────────────────────────────────────────────────────")
	}
}

func doAddToLibrary(scanner *bufio.Scanner) {
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

	if msg, ok := resp["message"]; ok {
		fmt.Println("✅", msg)
	} else {
		fmt.Println("❌", resp["error"])
	}
}

func doUpdateProgress(scanner *bufio.Scanner) {
	fmt.Print("Enter Manga ID: ")
	scanner.Scan()
	mangaID := strings.TrimSpace(scanner.Text())

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

	if msg, ok := resp["message"]; ok {
		fmt.Printf("✅ %s - Chapter %d\n", msg, chapter)
	} else {
		fmt.Println("❌", resp["error"])
	}
}

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

func (c *ChatClient) readLoop() {
	for {
		select {
		case <-c.done:
			return
		default:
			_, rawMsg, err := c.conn.ReadMessage()
			if err != nil {
				fmt.Println("\n❌ Disconnected from chat")
				c.close()
				return
			}

			var msg WSMessage
			if err := json.Unmarshal(rawMsg, &msg); err != nil {
				continue
			}

			switch msg.Type {
			case "chat":
				fmt.Printf("\n💬 [%s] %s: %s\n", msg.Timestamp, msg.Sender, msg.Content)
			case "dm":
				fmt.Printf("\n📩 [DM][%s] %s: %s\n", msg.Timestamp, msg.Sender, msg.Content)

				saveDM(LocalMessage{
					Type:      "dm",
					Sender:    msg.Sender,
					Recipient: c.username,
					Content:   msg.Content,
					Timestamp: msg.Timestamp,
				})
			case "system":
				fmt.Printf("\n● %s\n", msg.Content)
			case "join":
				fmt.Printf("\n👋 %s joined\n", msg.Sender)
			case "leave":
				fmt.Printf("\n👋 %s left\n", msg.Sender)
			case "list_rooms":
				fmt.Printf("\n📦 Rooms: %s\n", msg.Content)
			case "user_list":
				fmt.Printf("\n👥 Users: %s\n", msg.Content)
			case "error":
				fmt.Printf("\n❌ %s\n", msg.Content)
			}

			fmt.Printf("[%s] > ", c.room)
		}
	}
}

func (c *ChatClient) writeLoop() {
	for {
		select {
		case msg := <-c.send:
			data, _ := json.Marshal(msg)
			err := c.conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println("\n❌ Send failed")
				c.close()
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
		c.conn.Close()
	})
}

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
