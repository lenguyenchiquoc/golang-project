package main

import (
	"log"

	apiserver  "managahub/cmd/api-server"
	grpcserver "managahub/cmd/grpc-server"
	tcpserver  "managahub/cmd/tcp-server"
	udpserver  "managahub/cmd/udp-server"
	wsserver   "managahub/cmd/ws-server"
	"managahub/pkg/database"
)

func main() {
	db := database.InitDB("./mangahub.db")
	defer db.Close()
	jwtSecret := "ITITIU22134_LENGUYENCHIQUOC"
	tcpServer := tcpserver.Run()    
	udpServer := udpserver.Run()    
	grpcserver.Run(db, jwtSecret)             
	hub       := wsserver.Run()    

	log.Println("MangaHub starting...")
	log.Println("HTTP      : http://localhost:8080")
	log.Println("TCP Sync  : tcp://localhost:9090")
	log.Println("UDP Notify: udp://localhost:9091")
	log.Println("gRPC      : grpc://localhost:9092")
	log.Println("WebSocket : ws://localhost:8080/ws")

	apiserver.Run(tcpServer, udpServer, hub, db, jwtSecret)
}