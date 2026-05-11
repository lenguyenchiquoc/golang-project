package manga

import (
	"net/http"
	"strings"
	"time"

	"managahub/internal/udp"
	"managahub/pkg/models"

	"github.com/gin-gonic/gin"
)
 
type MangaHandler struct {
	Service *MangaService
	UDPServer *udp.NotificationServer
}
 
func NewMangaHandler(service *MangaService, udpServer *udp.NotificationServer) *MangaHandler {
	return &MangaHandler{
		Service:   service,
		UDPServer: udpServer,
	}
}

func (h *MangaHandler) CreateManga(c *gin.Context) {
	var req models.CreateMangaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	manga, err := h.Service.CreateManga(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.UDPServer != nil {
		h.UDPServer.Broadcast(udp.Notification{
			Type:       "new_manga", 
			MangaID:    manga.ID,
			MangaTitle: manga.Title,
			Message:    "Read it, guys1" + manga.Title,
			Timestamp:  time.Now().Unix(),
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Add manga successfully!",
		"manga":   manga,
	})
}

// GET /manga?query=...&genre=...&status=...&page=1&limit=10 
func (h *MangaHandler) SearchManga(c *gin.Context) {
    var req models.SearchMangaRequest
    c.ShouldBindQuery(&req)

    mangas, err := h.Service.SearchManga(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"results": mangas})
}
 
// GET /manga/:id 
func (h *MangaHandler) GetManga(c *gin.Context) {
    id := c.Param("id")
    manga, err := h.Service.GetMangaByID(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, manga)
}

func (h *MangaHandler) RateManga(c *gin.Context) {
	var req struct {
		MangaID string `json:"manga_id"`
		Rating  int32  `json:"rating"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")

	result, err := h.Service.RateManga(userID, req.MangaID, req.Rating)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, result)
}

// PUT /manga/:id
func (h *MangaHandler) UpdateManga(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateMangaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	manga, err := h.Service.UpdateManga(id, req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if h.UDPServer != nil {
		h.UDPServer.Broadcast(udp.Notification{
			Type:       "chapter_release",
			MangaID:    manga.ID,
			MangaTitle: manga.Title,
			Message:    "Manga updated: " + manga.Title,
			Timestamp:  time.Now().Unix(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Updated successfully",
		"manga":   manga,
	})
}