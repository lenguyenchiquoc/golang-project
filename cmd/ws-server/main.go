package wsserver
 
import "managahub/internal/websocket"
 
func Run() *wsocket.ChatHub {
	hub := wsocket.NewChatHub()
	go hub.Run()
	return hub
}
 