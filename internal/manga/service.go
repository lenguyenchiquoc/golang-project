package manga
 
import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"managahub/pkg/models"
)
 
type MangaService struct {
	DB *sql.DB
}
 
func NewMangaService(db *sql.DB) *MangaService {
	return &MangaService{DB: db}
}

func (m *MangaService) CreateManga(request models.CreateMangaRequest) (*models.Manga, error) {
	var count int

	query := "Select count(*) from manga where title = ? and author = ?"
	err := m.DB.QueryRow(query, request.Title, request.Author).Scan(&count)
	if err != nil {
		return nil, errors.New("Databases error")
	}
	if count > 0 {
		return nil, errors.New("Manga already exist")
	}

	if request.Status == "" {
		request.Status = "ongoing"
	}

	if request.Cover_url == "" {
		request.Cover_url = ""
	}

	mangaId := generateMangaID()
	query1 := `INSERT INTO manga (id,title, author, genres, status, total_chapters, description, cover_url)
	          VALUES (?,?, ?, ?, ?, ?, ?,?)`
	_, err = m.DB.Exec(query1,
		mangaId,
		request.Title,
		request.Author,
		request.Genres,
		request.Status,
		request.TotalChapters,
		request.Description,
		request.Cover_url,
	)
	if err != nil {
		return nil, errors.New("Can not add" + err.Error())
	}

	return &models.Manga{
		ID:            mangaId,
		Title:         request.Title,
		Author:        request.Author,
		Genres:        request.Genres,
		Status:        request.Status,
		TotalChapters: request.TotalChapters,
		Description:   request.Description,
		Cover_url: request.Cover_url,
	}, nil

}

func generateMangaID() string {
    return fmt.Sprintf("manga_%d", time.Now().UnixNano())
}

func (m *MangaService) SearchManga(request models.SearchMangaRequest) ([]models.Manga, error) {
	if request.Page <= 0 {
		request.Page = 1
	}
	if request.Limit <= 0 {
		request.Limit = 10
	}

	conditions := []string{}
	args := []interface{}{}

	if request.Query != "" {
		conditions = append(conditions, "(title LIKE ? OR author LIKE ?)")
		args = append(args, "%"+request.Query+"%", "%"+request.Query+"%")
	}
	if request.Genre != "" {
		conditions = append(conditions, "genres LIKE ?")
		args = append(args, "%"+request.Genre+"%")
	}
	if request.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, request.Status)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	offset := (request.Page - 1) * request.Limit
	query := fmt.Sprintf(`
		SELECT id, title, author, genres, status, total_chapters, description, cover_url
		FROM manga %s
		ORDER BY title ASC
		LIMIT ? OFFSET ?`, whereClause)

	args = append(args, request.Limit, offset)
	rows, err := m.DB.Query(query, args...)
	if err != nil {
		return nil, errors.New("Error search" + err.Error())
	}
	defer rows.Close()

	mangas := []models.Manga{}
	for rows.Next() {
		var manga models.Manga
		err := rows.Scan(
			&manga.ID,
			&manga.Title,
			&manga.Author,
			&manga.Genres,
			&manga.Status,
			&manga.TotalChapters,
			&manga.Description,
			&manga.Cover_url,
		)
		if err != nil {
			return nil, errors.New("Error" + err.Error())
		}
		mangas = append(mangas, manga)
	}

	return mangas, nil
}

func (s *MangaService) GetMangaByID(id string) (*models.Manga, error) {
	var m models.Manga
	query := `SELECT id, title, author, genres, status, total_chapters, description, cover_url
	          FROM manga WHERE id = ?`
 
	err := s.DB.QueryRow(query, id).Scan(
		&m.ID,
		&m.Title,
		&m.Author,
		&m.Genres,
		&m.Status,
		&m.TotalChapters,
		&m.Description,
		&m.Cover_url,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("Can not find Manga")
	}
	if err != nil {
		return nil, errors.New("Error database" + err.Error())
	}
 
	return &m, nil
}