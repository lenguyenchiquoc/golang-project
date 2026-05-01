package library

import (
	"database/sql"
	"errors"
	"time"

	"managahub/internal/tcp"
	"managahub/pkg/models"
)

type LibraryService struct {
	DB *sql.DB
	TCPServer *tcp.ProgressSyncServer
}

func NewLibraryService(db *sql.DB, tcpServer *tcp.ProgressSyncServer) *LibraryService {
    return &LibraryService{DB: db, TCPServer: tcpServer}
}

func (s *LibraryService) AddToLibrary(userID string, req models.AddToLibraryRequest) error {
	var count int
	err := s.DB.QueryRow(`SELECT COUNT(*) FROM manga WHERE id = ?`, req.MangaID).Scan(&count)
	if err != nil {
		return errors.New("Database is error " + err.Error())
	}
	if count == 0 {
		return errors.New("manga is not existed")
	}

	err = s.DB.QueryRow(
		`SELECT COUNT(*) FROM user_progress WHERE user_id = ? AND manga_id = ?`,
		userID, req.MangaID,
	).Scan(&count)
	if err != nil {
		return errors.New("Database is error" + err.Error())
	}
	if count > 0 {
		return errors.New("manga was already existed")
	}

	if !isValidStatus(req.Status) {
		return errors.New("Not acceptable status: reading, completed, plan_to_read")
	}

	query := `INSERT INTO user_progress (user_id, manga_id, current_chapter, status)
	          VALUES (?, ?, ?, ?)`
	_, err = s.DB.Exec(query, userID, req.MangaID, req.CurrentChapter, req.Status)
	if err != nil {
		return errors.New("Can not add " + err.Error())
	}

	return nil
}

func (s *LibraryService) GetLibrary(userID string) ([]models.LibraryItem, error) {
	query := `
		SELECT 
			up.manga_id,
			m.title,
			m.author,
			up.status,
			up.current_chapter,
			m.total_chapters,
			up.updated_at
		FROM user_progress up
		JOIN manga m ON up.manga_id = m.id
		WHERE up.user_id = ?
		ORDER BY up.updated_at DESC`

	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, errors.New("database is error: " + err.Error())
	}
	defer rows.Close()

	items := []models.LibraryItem{}
	for rows.Next() {
		var item models.LibraryItem
		err := rows.Scan(
			&item.MangaID,
			&item.Title,
			&item.Author,
			&item.Status,
			&item.CurrentChapter,
			&item.TotalChapters,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, errors.New("Error reading " + err.Error())
		}
		items = append(items, item)
	}

	return items, nil
}

func (s *LibraryService) UpdateProgress(userID string, req models.UpdateProgressRequest) error {
	var count int
	err := s.DB.QueryRow(
		`SELECT COUNT(*) FROM user_progress WHERE user_id = ? AND manga_id = ?`,
		userID, req.MangaID,
	).Scan(&count)
	if err != nil {
		return errors.New("Database is error: " + err.Error())
	}
	if count == 0 {
		return errors.New("Manga is not existed in library")
	}

	if req.CurrentChapter < 0 {
		return errors.New("No of chapter is not posible")
	}

	if req.Status != "" && !isValidStatus(req.Status) {
		return errors.New("Invalid status ")
	}

	if req.Status == "" {
		_, err = s.DB.Exec(
			`UPDATE user_progress SET current_chapter = ?, updated_at = CURRENT_TIMESTAMP
			 WHERE user_id = ? AND manga_id = ?`,
			req.CurrentChapter, userID, req.MangaID,
		)
	} else {
		_, err = s.DB.Exec(
			`UPDATE user_progress SET current_chapter = ?, status = ?, updated_at = CURRENT_TIMESTAMP
			 WHERE user_id = ? AND manga_id = ?`,
			req.CurrentChapter, req.Status, userID, req.MangaID,
		)
	}

	if err != nil {
		return errors.New("Can not update progress " + err.Error())
	}
	if s.TCPServer != nil {
    s.TCPServer.BroadcastProgress(tcp.ProgressUpdate{
        UserID:         userID,
        MangaID:        req.MangaID,
        CurrentChapter: req.CurrentChapter,
        Status:         req.Status,
        Timestamp:      time.Now().Unix(),
    })
}
	return nil
}

func isValidStatus(status string) bool {
	return status == "reading" || status == "completed" || status == "plan_to_read"
}
