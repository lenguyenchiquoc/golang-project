package grpcserver

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"

	pb "managahub/pkg/proto/managahub/pkg/proto"

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

	rows, err := s.DB.QueryContext(ctx, `
		SELECT 
			m.id, m.title, m.author, m.status, m.total_chapters, m.description, m.cover_url,
			g.id, g.name, m.average_rating, m.rating_count
		FROM manga m
		LEFT JOIN manga_genres mg ON m.id = mg.manga_id
		LEFT JOIN genres g ON mg.genre_id = g.id
		WHERE m.id = ?
	`, req.Id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var manga *pb.MangaResponse
	genres := []*pb.Genre{}

	for rows.Next() {
		var m pb.MangaResponse
		var gID sql.NullString
		var gName sql.NullString
		var avg float32
		var count int32
		err := rows.Scan(
			&m.Id,
			&m.Title,
			&m.Author,
			&m.Status,
			&m.TotalChapters,
			&m.Description,
			&m.CoverUrl,
			&gID,
			&gName,
			&avg,
			&count,
		)
		if err != nil {
			return nil, err
		}

		if manga == nil {
			manga = &m
			manga.AverageRating = avg
			manga.RatingCount = count
		}

		if gID.Valid {
			genres = append(genres, &pb.Genre{
				Id:   gID.String,
				Name: gName.String,
			})
		}
	}

	if manga == nil {
		return nil, errors.New("manga not found")
	}

	manga.Genres = genres
	return manga, nil
}

func (s *MangaGRPCServer) SearchManga(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	query := `
		SELECT 
			m.id, m.title, m.author, m.status, m.total_chapters, m.description, m.cover_url,
			g.id, g.name
		FROM manga m
		LEFT JOIN manga_genres mg ON m.id = mg.manga_id
		LEFT JOIN genres g ON mg.genre_id = g.id
		WHERE 1=1
	`

	args := []interface{}{}

	if req.Query != "" {
		query += " AND (m.title LIKE ? OR m.author LIKE ?)"
		args = append(args, "%"+req.Query+"%", "%"+req.Query+"%")
	}

	if len(req.Genre) > 0 {
		query += " AND g.name IN ("
		for i := range req.Genre {
			query += "?"
			if i < len(req.Genre)-1 {
				query += ","
			}
			args = append(args, req.Genre[i])
		}
		query += ")"
	}

	if req.Status != "" {
		query += " AND m.status = ?"
		args = append(args, req.Status)
	}

	query += " ORDER BY m.title ASC LIMIT ? OFFSET ?"

	offset := (req.Page - 1) * req.Limit
	args = append(args, req.Limit, offset)

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.New("database error: " + err.Error())
	}
	defer rows.Close()

	mangaMap := map[string]*pb.MangaResponse{}

	for rows.Next() {
		var (
			id, title, author, status, desc, cover string
			total                                  int32
			gID                                    sql.NullString
			gName                                  sql.NullString
		)

		err := rows.Scan(
			&id, &title, &author, &status,
			&total, &desc, &cover,
			&gID, &gName,
		)
		if err != nil {
			continue
		}

		if _, exists := mangaMap[id]; !exists {
			mangaMap[id] = &pb.MangaResponse{
				Id:            id,
				Title:         title,
				Author:        author,
				Status:        status,
				TotalChapters: total,
				Description:   desc,
				CoverUrl:      cover,
				Genres:        []*pb.Genre{},
			}
		}

		if gID.Valid {
			mangaMap[id].Genres = append(mangaMap[id].Genres, &pb.Genre{
				Id:   gID.String,
				Name: gName.String,
			})
		}
	}

	result := []*pb.MangaResponse{}
	for _, m := range mangaMap {
		result = append(result, m)
	}

	return &pb.SearchResponse{
		Mangas: result,
		Total:  int32(len(result)),
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

func (s *MangaGRPCServer) RateManga(ctx context.Context, req *pb.RatingRequest) (*pb.RatingResponse, error) {
	if req.Rating < 1 || req.Rating > 10 {
		return &pb.RatingResponse{
			Success: false,
			Message: "rating must be 1-10",
		}, nil
	}

	var exists int
	err := s.DB.QueryRow(`SELECT COUNT(*) FROM manga WHERE id = ?`, req.MangaId).Scan(&exists)
	if err != nil || exists == 0 {
		return &pb.RatingResponse{
			Success: false,
			Message: "manga not found",
		}, nil
	}

	_, err = s.DB.Exec(`
		INSERT INTO user_ratings (user_id, manga_id, rating)
		VALUES (?, ?, ?)
		ON CONFLICT(user_id, manga_id)
		DO UPDATE SET rating = excluded.rating
	`, req.UserId, req.MangaId, req.Rating)

	if err != nil {
		return &pb.RatingResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	_, err = s.DB.Exec(`
		UPDATE manga
		SET 
			average_rating = (SELECT AVG(rating) FROM user_ratings WHERE manga_id = ?),
			rating_count   = (SELECT COUNT(*) FROM user_ratings WHERE manga_id = ?)
		WHERE id = ?
	`, req.MangaId, req.MangaId, req.MangaId)

	if err != nil {
		return &pb.RatingResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	var avg float32
	var count int32

	s.DB.QueryRow(`
		SELECT average_rating, rating_count FROM manga WHERE id = ?
	`, req.MangaId).Scan(&avg, &count)

	return &pb.RatingResponse{
		Success: true,
		Message: "Rated successfully",
		Average: avg,
		Count:   count,
	}, nil
}
