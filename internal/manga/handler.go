package manga
 
import (
	"net/http"
 
	"github.com/gin-gonic/gin"

	grpcserver "managahub/internal/grpc"
	"managahub/pkg/models"
)
 
type MangaHandler struct {
	Service *MangaService
	GRPCClient *grpcserver.MangaGRPCClient
}
 
func NewMangaHandler(service *MangaService, grpcClient *grpcserver.MangaGRPCClient) *MangaHandler {
	return &MangaHandler{
		Service:    service,
		GRPCClient: grpcClient,
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
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "query params không hợp lệ",
		})
		return
	}
 
	resp, err := h.GRPCClient.SearchManga(
		req.Query,
		req.Genre,
		req.Status,
		int32(req.Page),
		int32(req.Limit),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":   resp.Total,
		"page":    req.Page,
		"limit":   req.Limit,
		"results": resp.Mangas,
	})
}
 
// GET /manga/:id 
func (h *MangaHandler) GetManga(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Lack id param"})
		return
	}
 
	manga, err := h.GRPCClient.GetManga(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, manga)
}