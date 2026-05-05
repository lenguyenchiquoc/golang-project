package grpcserver

import (
	"context"
	"log"
	"time"

	pb "managahub/pkg/proto/managahub/pkg/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type MangaGRPCClient struct {
	client pb.MangaServiceClient
	conn   *grpc.ClientConn
	getTokenFn     func() string 
	onAuthFailedFn func()
}

func NewMangaGRPCClient(address string, getToken func() string, onAuthFailed func()) *MangaGRPCClient {
	c := &MangaGRPCClient{
		getTokenFn:     getToken,
		onAuthFailedFn: onAuthFailed,
	}

	conn, err := grpc.Dial(address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(c.authInterceptor),
	)
	if err != nil {
		log.Fatal("Cannot connect to gRPC server:", err)
	}

	c.conn = conn
	c.client = pb.NewMangaServiceClient(conn)
	return c
}

func (c *MangaGRPCClient) Close() {
	c.conn.Close()
}

func (c *MangaGRPCClient) authInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	token := c.getTokenFn()
	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	err := invoker(ctx, method, req, reply, cc, opts...)

	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.Unauthenticated {
			log.Println("🚨 [gRPC] Token invalid or expired, triggering logout...")
			c.onAuthFailedFn()
		}
	}
	return err
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

func (c *MangaGRPCClient) RateManga(userID, mangaID string, rating int32) (*pb.RatingResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return c.client.RateManga(ctx, &pb.RatingRequest{
		UserId:  userID,
		MangaId: mangaID,
		Rating:  rating,
	})
}

func (c *MangaGRPCClient) GetClient() pb.MangaServiceClient {
    return c.client
}
