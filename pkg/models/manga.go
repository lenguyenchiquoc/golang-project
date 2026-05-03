package models

type Manga struct {
	ID            string `json:"id" binding:"required"`
	Title         string `json:"title" binding:"required"`
	Author        string `json:"author" binding:"required"`
	Genres        []Genre `json:"genres" binding:"required"`
	Status        string `json:"status" binding:"required"`
	TotalChapters int    `json:"total_chapter" binding:"required"`
	Description   string `json:"description" binding:"required"`
	Cover_url     string `json:"cover_url" binding:"required"`
}

type CreateMangaRequest struct {
	Title         string `json:"title" binding:"required"`
	Author        string `json:"author" binding:"required"`
	Genres        string `json:"genres"`
	Status        string `json:"status"`
	TotalChapters int    `json:"total_chapters"`
	Description   string `json:"description"`
	Cover_url     string `json:"cover_url"`
}
type SearchMangaRequest struct {
	Query  string `form:"query"`
	Genre  string `form:"genre"`
	Status string `form:"status"`
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
}

type Genre struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
