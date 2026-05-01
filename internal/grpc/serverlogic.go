package grpcserver

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"

	pb "managahub/pkg/proto"

	"google.golang.org/grpc"
)

type MangaGRPCServer struct {
	pb.UnimplementedMangaServiceServer
	DB *sql.DB
}

func NewMangaGRPCServer(db *sql.DB) *MangaGRPCServer {
	return &MangaGRPCServer{DB: db}
}

func (s *MangaGRPCServer) GetManga(ctx context.Context, req *pb.GetMangaRequest) (*pb.MangaResponse, error) {
	if req.Id == "" {
		return nil, errors.New("manga ID is required")
	}

	var m pb.MangaResponse
	query := `SELECT id, title, author, genres, status, total_chapters, description, cover_url
	          FROM manga WHERE id = ?`

	err := s.DB.QueryRowContext(ctx, query, req.Id).Scan(
		&m.Id, &m.Title, &m.Author, &m.Genres,
		&m.Status, &m.TotalChapters, &m.Description, &m.CoverUrl,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("manga not found")
	}
	if err != nil {
		return nil, errors.New("database error: " + err.Error())
	}

	return &m, nil
}

func (s *MangaGRPCServer) SearchManga(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	query := `SELECT id, title, author, genres, status, total_chapters, description, cover_url
	          FROM manga WHERE 1=1`
	args := []interface{}{}

	if req.Query != "" {
		query += " AND (title LIKE ? OR author LIKE ?)"
		args = append(args, "%"+req.Query+"%", "%"+req.Query+"%")
	}
	if req.Genre != "" {
		query += " AND genres LIKE ?"
		args = append(args, "%"+req.Genre+"%")
	}
	if req.Status != "" {
		query += " AND status = ?"
		args = append(args, req.Status)
	}

	offset := (req.Page - 1) * req.Limit
	query += " ORDER BY title ASC LIMIT ? OFFSET ?"
	args = append(args, req.Limit, offset)

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.New("database error: " + err.Error())
	}
	defer rows.Close()

	mangas := []*pb.MangaResponse{}
	for rows.Next() {
		var m pb.MangaResponse
		err := rows.Scan(
			&m.Id, &m.Title, &m.Author, &m.Genres,
			&m.Status, &m.TotalChapters, &m.Description, &m.CoverUrl,
		)
		if err != nil {
			continue
		}
		mangas = append(mangas, &m)
	}

	return &pb.SearchResponse{
		Mangas: mangas,
		Total:  int32(len(mangas)),
	}, nil
}

func (s *MangaGRPCServer) UpdateProgress(ctx context.Context, req *pb.ProgressRequest) (*pb.ProgressResponse, error) {
	if req.UserId == "" || req.MangaId == "" {
		return &pb.ProgressResponse{
			Success: false,
			Message: "user_id and manga_id are required",
		}, nil
	}

	if req.CurrentChapter < 0 {
		return &pb.ProgressResponse{
			Success: false,
			Message: "chapter cannot be negative",
		}, nil
	}

	var count int
	err := s.DB.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM user_progress WHERE user_id = ? AND manga_id = ?`,
		req.UserId, req.MangaId,
	).Scan(&count)
	if err != nil || count == 0 {
		return &pb.ProgressResponse{
			Success: false,
			Message: "manga not found in library",
		}, nil
	}

	_, err = s.DB.ExecContext(ctx,
		`UPDATE user_progress SET current_chapter = ?, status = ?, updated_at = CURRENT_TIMESTAMP
		 WHERE user_id = ? AND manga_id = ?`,
		req.CurrentChapter, req.Status, req.UserId, req.MangaId,
	)
	if err != nil {
		return &pb.ProgressResponse{
			Success: false,
			Message: "failed to update: " + err.Error(),
		}, nil
	}

	return &pb.ProgressResponse{
		Success: true,
		Message: "Progress updated successfully!",
	}, nil
}

func (s *MangaGRPCServer) Start(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Cannot start gRPC server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMangaServiceServer(grpcServer, s)

	log.Println("gRPC Server running at port", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("gRPC server error:", err)
	}
}
