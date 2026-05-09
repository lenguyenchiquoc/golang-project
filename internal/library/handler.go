package library

import (
	"net/http"
	"strings"

	"managahub/pkg/models"

	"github.com/gin-gonic/gin"
)

type LibraryHandler struct {
	Service    *LibraryService
}

func NewLibraryHandler(service *LibraryService) *LibraryHandler {
	return &LibraryHandler{
		Service:    service,
	}
}

// POST /users/library
func (h *LibraryHandler) AddToLibrary(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.AddToLibraryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.AddToLibrary(userID.(string), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Added to library successfully!",
		"manga_id": req.MangaID,
		"status":   req.Status,
	})
}

// GET /users/library
func (h *LibraryHandler) GetLibrary(c *gin.Context) {
	userID, _ := c.Get("user_id")

	items, err := h.Service.GetLibrary(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":   len(items),
		"library": items,
	})
}

// PUT /users/progress
func (h *LibraryHandler) UpdateProgress(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.UpdateProgress(userID.(string), req)
	if err != nil {
        if strings.Contains(err.Error(), "exist") || strings.Contains(err.Error(), "invalid") {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Error update/ pls check input or internal server is error"})
        }
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "status":          "success",
        "message":         "Progress updated successfully",
        "manga_id":        req.MangaID,
        "current_chapter": req.CurrentChapter,
    })
}