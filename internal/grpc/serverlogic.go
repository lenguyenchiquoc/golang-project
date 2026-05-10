package grpcserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	pb "managahub/pkg/proto/managahub/pkg/proto"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type MangaGRPCServer struct {
	pb.UnimplementedMangaServiceServer
	DB *sql.DB
	JWTSecret string
}

func NewMangaGRPCServer(db *sql.DB, secret string) *MangaGRPCServer {
    return &MangaGRPCServer{
        DB:        db,
        JWTSecret: secret,
    }
}

type Claims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
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

	// Bước 1: Lấy danh sách manga ID trước với LIMIT/OFFSET
	idQuery := `SELECT DISTINCT m.id FROM manga m
		LEFT JOIN manga_genres mg ON m.id = mg.manga_id
		LEFT JOIN genres g ON mg.genre_id = g.id
		WHERE 1=1`

	args := []interface{}{}

	if req.Query != "" {
		idQuery += " AND (m.title LIKE ? OR m.author LIKE ?)"
		args = append(args, "%"+req.Query+"%", "%"+req.Query+"%")
	}

	if len(req.Genre) > 0 {
		idQuery += " AND g.name IN ("
		for i := range req.Genre {
			idQuery += "?"
			if i < len(req.Genre)-1 {
				idQuery += ","
			}
			args = append(args, req.Genre[i])
		}
		idQuery += ")"
	}

	if req.Status != "" {
		idQuery += " AND m.status = ?"
		args = append(args, req.Status)
	}

	// Đếm total
	countQuery := "SELECT COUNT(*) FROM (" + idQuery + ") as sub"
	var total int32
	s.DB.QueryRowContext(ctx, countQuery, args...).Scan(&total)

	// Lấy ID với đúng LIMIT/OFFSET
	offset := (req.Page - 1) * req.Limit
	idQuery += " ORDER BY m.id ASC LIMIT ? OFFSET ?"
	idArgs := append(args, req.Limit, offset)

	idRows, err := s.DB.QueryContext(ctx, idQuery, idArgs...)
	if err != nil {
		return nil, errors.New("database error: " + err.Error())
	}

	var ids []string
	for idRows.Next() {
		var id string
		idRows.Scan(&id)
		ids = append(ids, id)
	}
	idRows.Close()

	if len(ids) == 0 {
		return &pb.SearchResponse{Mangas: []*pb.MangaResponse{}, Total: 0}, nil
	}

	// Bước 2: Lấy full data cho các ID đó
	placeholder := strings.Repeat("?,", len(ids))
	placeholder = placeholder[:len(placeholder)-1]

	dataQuery := fmt.Sprintf(`
		SELECT m.id, m.title, m.author, m.status, m.total_chapters, m.description, m.cover_url,
			g.id, g.name
		FROM manga m
		LEFT JOIN manga_genres mg ON m.id = mg.manga_id
		LEFT JOIN genres g ON mg.genre_id = g.id
		WHERE m.id IN (%s)
		ORDER BY m.title ASC`, placeholder)

	dataArgs := make([]interface{}, len(ids))
	for i, id := range ids {
		dataArgs[i] = id
	}

	rows, err := s.DB.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return nil, errors.New("database error: " + err.Error())
	}
	defer rows.Close()

	mangaMap := map[string]*pb.MangaResponse{}
	orderMap := []string{} // giữ thứ tự

	for rows.Next() {
		var (
			id, title, author, status, desc, cover string
			total_ch                               int32
			gID                                    sql.NullString
			gName                                  sql.NullString
		)
		err := rows.Scan(&id, &title, &author, &status, &total_ch, &desc, &cover, &gID, &gName)
		if err != nil {
			continue
		}

		if _, exists := mangaMap[id]; !exists {
			mangaMap[id] = &pb.MangaResponse{
				Id:            id,
				Title:         title,
				Author:        author,
				Status:        status,
				TotalChapters: total_ch,
				Description:   desc,
				CoverUrl:      cover,
				Genres:        []*pb.Genre{},
			}
			orderMap = append(orderMap, id)
		}

		if gID.Valid {
			mangaMap[id].Genres = append(mangaMap[id].Genres, &pb.Genre{
				Id:   gID.String,
				Name: gName.String,
			})
		}
	}

	result := make([]*pb.MangaResponse, 0, len(orderMap))
	for _, id := range orderMap {
		result = append(result, mangaMap[id])
	}

	return &pb.SearchResponse{
		Mangas: result,
		Total:  total,
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
    grpcServer := grpc.NewServer(
        grpc.UnaryInterceptor(s.authInterceptor),
    )

    pb.RegisterMangaServiceServer(grpcServer, s)

    log.Println("🚀 gRPC Server running at port", port)
    if err := grpcServer.Serve(listener); err != nil {
        log.Fatal("gRPC server error:", err)
    }
}

func (s *MangaGRPCServer) authInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Yêu cầu không có metadata")
	}

	authHeader := md["authorization"]
	// if len(authHeader) == 0 {
	// 	return nil, status.Error(codes.Unauthenticated, "Bạn chưa đăng nhập (Thiếu Token)")
	// }

	tokenStr := strings.TrimPrefix(authHeader[0], "Bearer ")
	claims, err := s.JWTParser(tokenStr) 
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Token is expired or invalid")
	}
	newCtx := context.WithValue(ctx, "user_id", claims.UserID)
	return handler(newCtx, req)
}

func (s *MangaGRPCServer) JWTParser(tokenStr string) (*Claims, error) {

    key := []byte(s.JWTSecret)

    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return key, nil
    })

    if err != nil {
        return nil, err 
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}



func (s *MangaGRPCServer) validateToken(token string) (bool, string) {
	// Logic parse JWT và check Expired ở đây
	// Trả về (true, userID) nếu hợp lệ, (false, "") nếu hết hạn
	return true, "some-user-id"
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
