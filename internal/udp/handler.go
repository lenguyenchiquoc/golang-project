package udp

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	Server *NotificationServer
}

func NewNotificationHandler(server *NotificationServer) *NotificationHandler {
	return &NotificationHandler{Server: server}
}

func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var req struct {
		MangaID    string `json:"manga_id"    binding:"required"`
		MangaTitle string `json:"manga_title" binding:"required"`
		Message    string `json:"message"     binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification := Notification{
		Type:       "chapter_release",
		MangaID:    req.MangaID,
		MangaTitle: req.MangaTitle,
		Message:    req.Message,
		Timestamp:  time.Now().Unix(),
	}

	h.Server.Broadcast(notification)

	c.JSON(http.StatusOK, gin.H{
		"message":       "Đã gửi thông báo!",
		"clients_count": h.Server.GetClientCount(),
	})
}
