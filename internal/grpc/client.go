package grpcserver

import (
	"context"
	"log"
	"time"

	pb "managahub/pkg/proto/managahub/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MangaGRPCClient struct {
	client pb.MangaServiceClient
	conn   *grpc.ClientConn
}

func NewMangaGRPCClient(address string) *MangaGRPCClient {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Cannot connect to gRPC server:", err)
	}

	return &MangaGRPCClient{
		client: pb.NewMangaServiceClient(conn),
		conn:   conn,
	}
}

func (c *MangaGRPCClient) Close() {
	c.conn.Close()
}

func (c *MangaGRPCClient) GetManga(id string) (*pb.MangaResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.client.GetManga(ctx, &pb.GetMangaRequest{Id: id})
}

func (c *MangaGRPCClient) SearchManga(query string, genres []string, status string, page, limit int32) (*pb.SearchResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.client.SearchManga(ctx, &pb.SearchRequest{
		Query:  query,
		Genre:  genres,
		Status: status,
		Page:   page,
		Limit:  limit,
	})
}

func (c *MangaGRPCClient) UpdateProgress(userID, mangaID string, chapter int32, status string) (*pb.ProgressResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.client.UpdateProgress(ctx, &pb.ProgressRequest{
		UserId:         userID,
		MangaId:        mangaID,
		CurrentChapter: chapter,
		Status:         status,
	})
}

func (c *MangaGRPCClient) RateManga(userID, mangaID string, rating int32) (*pb.RatingResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.client.RateManga(ctx, &pb.RatingRequest{
		UserId:  userID,
		MangaId: mangaID,
		Rating:  rating,
	})
}