package tcpserver

import (
	"log"
	"net"

	"managahub/internal/tcp"
)

func Run() *tcp.ProgressSyncServer {
	server := tcp.NewProgressSyncServer("9090")

	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatal("can not run TCP server:", err)
	}

	log.Println("✓ TCP Sync Server is running at port 9090")

	go server.HandleMessages()

	go func() {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Errror accept connection:", err)
				continue
			}
			go server.HandleClient(conn)
		}
	}()

	return server
}