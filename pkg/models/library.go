package models

import "time"

type UserProgress struct {
	UserID         string    `json:"user_id"`
	MangaID        string    `json:"manga_id"`
	CurrentChapter int       `json:"current_chapter"`
	Status         string    `json:"status"` 
	UpdatedAt      time.Time `json:"updated_at"`
}

// Request add
type AddToLibraryRequest struct {
	MangaID        string `json:"manga_id"        binding:"required"`
	Status         string `json:"status"          binding:"required"`
	CurrentChapter int    `json:"current_chapter"`
}

// Update reading progress
type UpdateProgressRequest struct {
	MangaID        string `json:"manga_id"        binding:"required"`
	CurrentChapter int    `json:"current_chapter" binding:"required"`
	Status         string `json:"status"`
}

// Response 
type LibraryItem struct {
	MangaID        string    `json:"manga_id"`
	Title          string    `json:"title"`
	Author         string    `json:"author"`
	Status         string    `json:"status"`
	CurrentChapter int       `json:"current_chapter"`
	TotalChapters  int       `json:"total_chapters"`
	UpdatedAt      time.Time `json:"updated_at"`
}
