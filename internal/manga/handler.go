package manga
 
import (
	"net/http"
 
	"github.com/gin-gonic/gin"
	"managahub/pkg/models"
)
 
type MangaHandler struct {
	Service *MangaService
}
 
func NewMangaHandler(service *MangaService) *MangaHandler {
	return &MangaHandler{
		Service:    service,
	}
}

func (h *MangaHandler) CreateManga(c *gin.Context) {
	var req models.CreateMangaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
 
	manga, err := h.Service.CreateManga(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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