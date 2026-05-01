package grpcserver

import (
	"database/sql"

	grpc "managahub/internal/grpc"
)

func Run(db *sql.DB) {
	server := grpc.NewMangaGRPCServer(db)
	go server.Start("9092")
}