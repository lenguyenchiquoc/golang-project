package library

import (
	"net/http"

	"github.com/gin-gonic/gin"

	grpcserver "managahub/internal/grpc"
	"managahub/pkg/models"
)

type LibraryHandler struct {
	Service    *LibraryService
	GRPCClient *grpcserver.MangaGRPCClient
}

func NewLibraryHandler(service *LibraryService, grpcClient *grpcserver.MangaGRPCClient) *LibraryHandler {
	return &LibraryHandler{
		Service:    service,
		GRPCClient: grpcClient,
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

// PUT /users/progress → gọi qua gRPC
func (h *LibraryHandler) UpdateProgress(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.GRPCClient.UpdateProgress(
		userID.(string),
		req.MangaID,
		int32(req.CurrentChapter),
		req.Status,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !resp.Success {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         resp.Message,
		"manga_id":        req.MangaID,
		"current_chapter": req.CurrentChapter,
	})
}