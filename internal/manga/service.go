package manga

import (
	"database/sql"
	"errors"
	"fmt"
	"managahub/pkg/models"
	"strings"
	"time"

	"github.com/google/uuid"
)

type MangaService struct {
	DB *sql.DB
}

func NewMangaService(db *sql.DB) *MangaService {
	return &MangaService{DB: db}
}

func (m *MangaService) CreateManga(req models.CreateMangaRequest) (*models.Manga, error) {
	tx, err := m.DB.Begin()
	if err != nil {
		return nil, err
	}

	var count int
	err = tx.QueryRow(`SELECT COUNT(*) FROM manga WHERE title = ? AND author = ?`,
		req.Title, req.Author).Scan(&count)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if count > 0 {
		tx.Rollback()
		return nil, errors.New("manga already exists")
	}

	if req.Status == "" {
		req.Status = "ongoing"
	}

	mangaID := generateMangaID()

	_, err = tx.Exec(`
		INSERT INTO manga (id, title, author, status, total_chapters, description, cover_url)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		mangaID,
		req.Title,
		req.Author,
		req.Status,
		req.TotalChapters,
		req.Description,
		req.Cover_url,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	genres := []models.Genre{}
	seen := map[string]bool{}

	for _, g := range req.Genres {
		g = strings.ToLower(strings.TrimSpace(g))
		if g == "" {
			continue
		}

		if seen[g] {
			continue
		}
		seen[g] = true

		var genreID string

		err := tx.QueryRow(`SELECT id FROM genres WHERE name = ?`, g).Scan(&genreID)

		if err == sql.ErrNoRows {
			genreID = uuid.New().String()
			_, err = tx.Exec(`INSERT INTO genres (id, name) VALUES (?, ?)`, genreID, g)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		} else if err != nil {
			tx.Rollback()
			return nil, err
		}

		_, err = tx.Exec(`
			INSERT INTO manga_genres (manga_id, genre_id)
			VALUES (?, ?)`,
			mangaID, genreID,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		genres = append(genres, models.Genre{
			ID:   genreID,
			Name: g,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Manga{
		ID:            mangaID,
		Title:         req.Title,
		Author:        req.Author,
		Status:        req.Status,
		TotalChapters: req.TotalChapters,
		Description:   req.Description,
		Cover_url:     req.Cover_url,
		Genres:        genres,
	}, nil
}

func generateMangaID() string {
	return fmt.Sprintf("manga_%d", time.Now().UnixNano())
}

func (m *MangaService) SearchManga(req models.SearchMangaRequest) ([]models.Manga, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	conditions := []string{}
	args := []interface{}{}

	if req.Query != "" {
		conditions = append(conditions, "(m.title LIKE ? OR m.author LIKE ?)")
		args = append(args, "%"+req.Query+"%", "%"+req.Query+"%")
	}

	if req.Genre != "" {
		conditions = append(conditions, "g.name = ?")
		args = append(args, req.Genre)
	}

	if req.Status != "" {
		conditions = append(conditions, "m.status = ?")
		args = append(args, req.Status)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	offset := (req.Page - 1) * req.Limit

	query := fmt.Sprintf(`
		SELECT DISTINCT m.id, m.title, m.author, m.status, m.total_chapters, m.description, m.cover_url, m.average_rating, m.rating_count
		FROM manga m
		LEFT JOIN manga_genres mg ON m.id = mg.manga_id
		LEFT JOIN genres g ON mg.genre_id = g.id
		%s
		ORDER BY m.title ASC
		LIMIT ? OFFSET ?`, whereClause)

	args = append(args, req.Limit, offset)

	rows, err := m.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mangas []models.Manga

	for rows.Next() {
		var manga models.Manga
		err := rows.Scan(
			&manga.ID,
			&manga.Title,
			&manga.Author,
			&manga.Status,
			&manga.TotalChapters,
			&manga.Description,
			&manga.Cover_url,
			&manga.AverageRating,
			&manga.RatingCount,
		)
		if err != nil {
			return nil, err
		}

		manga.Genres = m.getGenresByMangaID(manga.ID)

		mangas = append(mangas, manga)
	}

	return mangas, nil
}
func (s *MangaService) GetMangaByID(id string) (*models.Manga, error) {
	var m models.Manga

	err := s.DB.QueryRow(`
		SELECT id, title, author, status, total_chapters, description, cover_url, average_rating, rating_count
		FROM manga WHERE id = ?`, id,
	).Scan(
		&m.ID,
		&m.Title,
		&m.Author,
		&m.Status,
		&m.TotalChapters,
		&m.Description,
		&m.Cover_url,
		&m.AverageRating,
		&m.RatingCount,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("not found")
	}
	if err != nil {
		return nil, err
	}

	m.Genres = s.getGenresByMangaID(id)

	return &m, nil
}

func (s *MangaService) getGenresByMangaID(mangaID string) []models.Genre {
	rows, err := s.DB.Query(`
		SELECT g.id, g.name
		FROM genres g
		JOIN manga_genres mg ON g.id = mg.genre_id
		WHERE mg.manga_id = ?`, mangaID)

	if err != nil {
		return nil
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var g models.Genre
		rows.Scan(&g.ID, &g.Name)
		genres = append(genres, g)
	}

	return genres
}

func (m *MangaService) RateManga(userID, mangaID string, rating int32) (*models.RatingResponse, error) {
	if userID == "" || mangaID == "" {
		return nil, errors.New("user_id and manga_id are required")
	}

	if rating < 1 || rating > 10 {
		return nil, errors.New("rating must be between 1 and 10")
	}

	var exists int
	err := m.DB.QueryRow(`SELECT COUNT(*) FROM manga WHERE id = ?`, mangaID).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists == 0 {
		return nil, errors.New("manga not found")
	}

	_, err = m.DB.Exec(`
		INSERT INTO user_ratings (user_id, manga_id, rating)
		VALUES (?, ?, ?)
		ON CONFLICT(user_id, manga_id)
		DO UPDATE SET rating = excluded.rating
	`, userID, mangaID, rating)

	if err != nil {
		return nil, err
	}

	_, err = m.DB.Exec(`
		UPDATE manga
		SET 
			average_rating = (
				SELECT AVG(rating) FROM user_ratings WHERE manga_id = ?
			),
			rating_count = (
				SELECT COUNT(*) FROM user_ratings WHERE manga_id = ?
			)
		WHERE id = ?
	`, mangaID, mangaID, mangaID)

	if err != nil {
		return nil, err
	}

	var avg float32
	var count int32

	err = m.DB.QueryRow(`
		SELECT average_rating, rating_count
		FROM manga WHERE id = ?
	`, mangaID).Scan(&avg, &count)

	if err != nil {
		return nil, err
	}

	return &models.RatingResponse{
		Success: true,
		Message: "Rated successfully",
		Average: avg,
		Count:   count,
	}, nil
}
