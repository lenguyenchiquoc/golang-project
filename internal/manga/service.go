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
	for _, g := range req.Genres {
		var genreID string

		err := tx.QueryRow(`SELECT id FROM genres WHERE name = ?`, g).Scan(&genreID)

		if err == sql.ErrNoRows {
			genreID = uuid.New().String()
			_, err = tx.Exec(`INSERT INTO genres (id, name) VALUES (?, ?)`, genreID, g)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
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
	}

	tx.Commit()

	return &models.Manga{
		ID:     mangaID,
		Title:  req.Title,
		Author: req.Author,
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
		SELECT DISTINCT m.id, m.title, m.author, m.status, m.total_chapters, m.description, m.cover_url
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
		SELECT id, title, author, status, total_chapters, description, cover_url
		FROM manga WHERE id = ?`, id,
	).Scan(
		&m.ID,
		&m.Title,
		&m.Author,
		&m.Status,
		&m.TotalChapters,
		&m.Description,
		&m.Cover_url,
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
