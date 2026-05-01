package udpserver

import (
	"log"
	"net"

	"managahub/internal/udp"
)

func Run() *udp.NotificationServer {
	server := udp.NewNotificationServer("9091")

	addr, err := net.ResolveUDPAddr("udp", ":9091")
	if err != nil {
		log.Fatal("Cannot resolve UDP addr:", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal("Cannot start UDP server:", err)
	}

	server.SetConn(conn)

	log.Println("UDP Notification Server running at port 9091")

	buf := make([]byte, 1024)
	go func() {
		for {
			n, clientAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				log.Println("Error reading UDP:", err)
				continue
			}
			go server.HandlePacket(conn, buf[:n], clientAddr)
		}
	}()

	return server
}
