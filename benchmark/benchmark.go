package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	grpcserver "managahub/internal/grpc"

	"github.com/gorilla/websocket"
)

const (
	HTTP_SERVER = "http://localhost:8080"
	TCP_SERVER  = "localhost:9090"
	UDP_SERVER  = "localhost:9091"
	WS_SERVER   = "ws://localhost:8080/ws"
	NUM_CLIENTS = 50
)

// ====================== STRUCTS ======================

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
	Sender    string `json:"sender,omitempty"`
	Content   string `json:"content,omitempty"`
	Room      string `json:"room,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

type Result struct {
	ClientIdx int
	Action    string
	Success   bool
	Duration  time.Duration
	ErrMsg    string
}

// ====================== BENCH CLIENT ======================

type BenchClient struct {
	idx      int
	username string
	token    string
	userID   string

	tcpConn net.Conn
	udpConn *net.UDPConn
	wsConn  *websocket.Conn
	wsMu    sync.Mutex

	grpc *grpcserver.MangaGRPCClient

	// ws round-trip: marker -> signal channel
	wsPending   map[string]chan struct{}
	wsPendingMu sync.Mutex
	wsDone      chan struct{}
}

var (
	results []Result
	resMu   sync.Mutex
)

func addResult(r Result) {
	resMu.Lock()
	results = append(results, r)
	resMu.Unlock()
}

// ====================== HTTP ======================

func (c *BenchClient) httpRequest(method, path string, body interface{}) (map[string]interface{}, error) {
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
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
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

// ====================== AUTH ======================

func (c *BenchClient) register() {
	start := time.Now()
	resp, err := c.httpRequest("POST", "/auth/register", map[string]string{
		"username":   c.username,
		"email":      fmt.Sprintf("%s@bench.com", c.username),
		"password":   "password123",
		"RePassword": "password123",
	})
	elapsed := time.Since(start)
	addResult(Result{c.idx, "register", err == nil && resp != nil, elapsed, fmt.Sprintf("%v", err)})
}

func (c *BenchClient) login() bool {
	start := time.Now()
	resp, err := c.httpRequest("POST", "/auth/login", map[string]string{
		"username": c.username,
		"password": "password123",
	})
	elapsed := time.Since(start)
	if err != nil || resp["token"] == nil {
		addResult(Result{c.idx, "login", false, elapsed, fmt.Sprintf("err=%v resp=%v", err, resp)})
		return false
	}
	c.token = resp["token"].(string)
	c.userID = resp["user_id"].(string)
	addResult(Result{c.idx, "login", true, elapsed, ""})
	return true
}

func (c *BenchClient) logout() {
	c.httpRequest("POST", "/auth/logout", nil)
}

// ====================== TCP ======================

func (c *BenchClient) connectTCP() {
	start := time.Now()
	conn, err := net.Dial("tcp", TCP_SERVER)
	elapsed := time.Since(start)
	if err != nil {
		addResult(Result{c.idx, "tcp_connect", false, elapsed, err.Error()})
		return
	}
	c.tcpConn = conn

	auth := map[string]interface{}{
		"type": "auth",
		"payload": map[string]string{
			"user_id":  c.userID,
			"username": c.username,
		},
	}
	data, _ := json.Marshal(auth)
	conn.Write(append(data, '\n'))
	addResult(Result{c.idx, "tcp_connect", true, elapsed, ""})

	go func() {
		buf := make([]byte, 4096)
		for {
			conn.SetReadDeadline(time.Now().Add(30 * time.Second))
			_, err := conn.Read(buf)
			if err != nil {
				break
			}
		}
	}()
}

// ====================== UDP ======================

func (c *BenchClient) connectUDP() {
	start := time.Now()
	serverAddr, err := net.ResolveUDPAddr("udp", UDP_SERVER)
	if err != nil {
		addResult(Result{c.idx, "udp_connect", false, time.Since(start), err.Error()})
		return
	}
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		addResult(Result{c.idx, "udp_connect", false, time.Since(start), err.Error()})
		return
	}
	c.udpConn = conn

	packet := map[string]string{"type": "register", "user_id": c.userID}
	data, _ := json.Marshal(packet)
	conn.Write(data)
	addResult(Result{c.idx, "udp_connect", true, time.Since(start), ""})

	go func() {
		buf := make([]byte, 1024)
		for {
			conn.SetReadDeadline(time.Now().Add(30 * time.Second))
			n, err := conn.Read(buf)
			if err != nil {
				break
			}
			var notif UDPNotification
			json.Unmarshal(buf[:n], &notif)
		}
	}()
}

// ====================== GRPC ======================

func (c *BenchClient) initGRPC() {
	getToken := func() string { return c.token }
	onAuthFailed := func() { c.token = "" }
	c.grpc = grpcserver.NewMangaGRPCClient("localhost:9092", getToken, onAuthFailed)
}

func (c *BenchClient) searchMangaGRPC() {
	queries := []string{"", "Naruto", "One Piece", "Bleach", "action"}
	q := queries[rand.Intn(len(queries))]
	start := time.Now()
	resp, err := c.grpc.SearchManga(q, nil, "", 1, 10)
	elapsed := time.Since(start)
	addResult(Result{c.idx, "grpc_search", err == nil && resp != nil, elapsed, fmt.Sprintf("%v", err)})
}

func (c *BenchClient) getMangaGRPC() {
	resp, err := c.grpc.SearchManga("", nil, "", 1, 10)
	if err != nil || len(resp.Mangas) == 0 {
		addResult(Result{c.idx, "grpc_get_manga", false, 0, "no manga found"})
		return
	}
	id := resp.Mangas[rand.Intn(len(resp.Mangas))].Id
	start := time.Now()
	_, err = c.grpc.GetManga(id)
	elapsed := time.Since(start)
	addResult(Result{c.idx, "grpc_get_manga", err == nil, elapsed, fmt.Sprintf("%v", err)})
}

// ====================== HTTP LIBRARY ======================

func (c *BenchClient) getLibrary() {
	start := time.Now()
	resp, err := c.httpRequest("GET", "/users/library", nil)
	elapsed := time.Since(start)
	addResult(Result{c.idx, "get_library", err == nil && resp != nil, elapsed, fmt.Sprintf("%v", err)})
}

func (c *BenchClient) wsReadLoop() {
	defer func() {
		c.wsPendingMu.Lock()
		for marker, ch := range c.wsPending {
			close(ch)
			delete(c.wsPending, marker)
		}
		c.wsPendingMu.Unlock()
	}()

	for {
		select {
		case <-c.wsDone:
			return
		default:
		}

		c.wsConn.SetReadDeadline(time.Now().Add(30 * time.Second))
		_, raw, err := c.wsConn.ReadMessage()
		if err != nil {
			return
		}

		var msg WSMessage
		if json.Unmarshal(raw, &msg) != nil {
			continue
		}
		if msg.Type == "chat" && msg.Sender == c.username {
			c.wsPendingMu.Lock()
			ch, ok := c.wsPending[msg.Content]
			if ok {
				close(ch)
				delete(c.wsPending, msg.Content)
			}
			c.wsPendingMu.Unlock()
		}
	}
}

func (c *BenchClient) connectWSRoom(room string) bool {
	start := time.Now()
	url := fmt.Sprintf("%s?username=%s&userid=%s&room=%s", WS_SERVER, c.username, c.userID, room)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	elapsed := time.Since(start)
	if err != nil {
		addResult(Result{c.idx, "ws_connect", false, elapsed, err.Error()})
		return false
	}
	c.wsConn = conn
	c.wsPending = make(map[string]chan struct{})
	c.wsDone = make(chan struct{})
	addResult(Result{c.idx, "ws_connect", true, elapsed, ""})

	go c.wsReadLoop()
	return true
}

func (c *BenchClient) sendRoomMessageRoundTrip() {
	if c.wsConn == nil {
		return
	}

	marker := fmt.Sprintf("bench_%d_%d", c.idx, time.Now().UnixNano())

	waitCh := make(chan struct{})
	c.wsPendingMu.Lock()
	c.wsPending[marker] = waitCh
	c.wsPendingMu.Unlock()

	msg := WSMessage{Type: "chat", Content: marker}
	data, _ := json.Marshal(msg)

	start := time.Now()
	c.wsMu.Lock()
	err := c.wsConn.WriteMessage(websocket.TextMessage, data)
	c.wsMu.Unlock()

	if err != nil {
		c.wsPendingMu.Lock()
		delete(c.wsPending, marker)
		c.wsPendingMu.Unlock()
		addResult(Result{c.idx, "ws_roundtrip", false, time.Since(start), err.Error()})
		return
	}

	select {
	case <-waitCh:
		addResult(Result{c.idx, "ws_roundtrip", true, time.Since(start), ""})
	case <-time.After(15 * time.Second):
		c.wsPendingMu.Lock()
		delete(c.wsPending, marker)
		c.wsPendingMu.Unlock()
		addResult(Result{c.idx, "ws_roundtrip", false, time.Since(start), "timeout"})
	}
}


func (c *BenchClient) closeAll() {
	if c.wsDone != nil {
		close(c.wsDone)
	}
	if c.wsConn != nil {
		c.wsConn.Close()
	}
	if c.tcpConn != nil {
		c.tcpConn.Close()
	}
	if c.udpConn != nil {
		packet := map[string]string{"type": "unregister", "user_id": c.userID}
		data, _ := json.Marshal(packet)
		c.udpConn.Write(data)
		c.udpConn.Close()
	}
	if c.grpc != nil {
		c.grpc.Close()
	}
	c.logout()
}

// ====================== RUN CLIENT ======================

func runClient(idx int, wg *sync.WaitGroup, ready *sync.WaitGroup, startGun *sync.WaitGroup) {
	defer wg.Done()

	c := &BenchClient{
		idx:      idx,
		username: fmt.Sprintf("bench_%d", idx),
	}

	c.register()
	if !c.login() {
		ready.Done()
		return
	}
	c.initGRPC()
	c.connectTCP()
	c.connectUDP()

	ready.Done()
	startGun.Wait()

	defer c.closeAll()

	var actionWg sync.WaitGroup

	// gRPC search x3
	actionWg.Add(1)
	go func() {
		defer actionWg.Done()
		for i := 0; i < 3; i++ {
			c.searchMangaGRPC()
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		}
	}()

	actionWg.Add(1)
	go func() {
		defer actionWg.Done()
		c.getMangaGRPC()
	}()

	// HTTP get library
	actionWg.Add(1)
	go func() {
		defer actionWg.Done()
		c.getLibrary()
	}()

	actionWg.Add(1)
	go func() {
		defer actionWg.Done()
		if c.connectWSRoom("general") {
			for i := 0; i < 3; i++ {
				c.sendRoomMessageRoundTrip()
				time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
			}
		}
	}()

	actionWg.Wait()
}

// ====================== REPORT ======================

func printReport(totalElapsed time.Duration) {
	type Stats struct {
		success, failed int
		total, min, max time.Duration
	}

	statsMap := map[string]*Stats{}
	for _, r := range results {
		if _, ok := statsMap[r.Action]; !ok {
			statsMap[r.Action] = &Stats{min: 999 * time.Second}
		}
		s := statsMap[r.Action]
		if r.Success {
			s.success++
		} else {
			s.failed++
		}
		s.total += r.Duration
		if r.Duration > 0 && r.Duration < s.min {
			s.min = r.Duration
		}
		if r.Duration > s.max {
			s.max = r.Duration
		}
	}

	fmt.Println("\n╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                        BENCHMARK REPORT                             ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
	fmt.Printf("%-22s %6s %6s %10s %10s %10s\n", "Action", "OK", "FAIL", "Avg", "Min", "Max")
	fmt.Println(strings.Repeat("─", 70))

	actions := []string{
		"register", "login",
		"tcp_connect", "udp_connect",
		"grpc_search", "grpc_get_manga",
		"get_library",
		"ws_connect", "ws_roundtrip",
	}

	allOK := true
	for _, action := range actions {
		s, ok := statsMap[action]
		if !ok {
			continue
		}
		total := s.success + s.failed
		var avg time.Duration
		if total > 0 {
			avg = s.total / time.Duration(total)
		}
		icon := "✅"
		if s.failed > 0 {
			icon = "❌"
			allOK = false
		}
		fmt.Printf("%s %-20s %6d %6d %10v %10v %10v\n",
			icon, action,
			s.success, s.failed,
			avg.Round(time.Millisecond),
			s.min.Round(time.Millisecond),
			s.max.Round(time.Millisecond),
		)
	}

	fmt.Println(strings.Repeat("─", 70))
	fmt.Printf("Clients      : %d\n", NUM_CLIENTS)
	fmt.Printf("Total actions: %d\n", len(results))
	fmt.Printf("Total time   : %v\n", totalElapsed.Round(time.Millisecond))

	if allOK {
		fmt.Println("\n🎉 Server handled all 50 clients successfully!")
	} else {
		fmt.Println("\n⚠️  Some actions failed!")
		fmt.Println("\nFailed details:")
		printed := 0
		for _, r := range results {
			if !r.Success && r.ErrMsg != "<nil>" && r.ErrMsg != "" {
				fmt.Printf("  [client_%d] %s: %s\n", r.ClientIdx, r.Action, r.ErrMsg)
				printed++
				if printed >= 10 {
					fmt.Println("  ... (truncated)")
					break
				}
			}
		}
	}
}

// ====================== MAIN ======================

func main() {
	fmt.Printf("🚀 Benchmark: %d concurrent clients\n", NUM_CLIENTS)
	fmt.Println("   Protocols: HTTP + gRPC + TCP + UDP + WebSocket (round-trip, same room)")
	fmt.Println()

	var wg, ready sync.WaitGroup
	var startGun sync.WaitGroup

	ready.Add(NUM_CLIENTS)
	startGun.Add(1)

	for i := 0; i < NUM_CLIENTS; i++ {
		wg.Add(1)
		go runClient(i, &wg, &ready, &startGun)
		time.Sleep(30 * time.Millisecond)
	}

	fmt.Println("⏳ Waiting for all clients to setup...")
	ready.Wait()

	fmt.Printf("🔥 All %d clients ready! Firing...\n", NUM_CLIENTS)
	start := time.Now()
	startGun.Done()

	wg.Wait()
	printReport(time.Since(start))
}
