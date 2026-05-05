package grpcserver

import (
	"database/sql"

	grpc "managahub/internal/grpc"
)

func Run(db *sql.DB, jwt string) {
	server := grpc.NewMangaGRPCServer(db,jwt)
	go server.Start("9092")
}