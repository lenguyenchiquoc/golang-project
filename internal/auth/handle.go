package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"managahub/pkg/models"
)

type AuthHandler struct {
	Service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request errors",
		})
		return
	}
	resp, err := h.Service.login(req)
	if err != nil {
		switch err.Error() {
		case "1":
			c.JSON(http.StatusNotFound, gin.H{"error": "Doesnt exist"})
		case "2":
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong password"})
		case "already_logged_in":
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Account already logged in on another device",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Login success!",
		"token":    resp.Token,
		"username": resp.Username,
		"user_id":  resp.UserID,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID := c.GetString("user_id") 

	Logout(userID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out",
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := h.Service.register(req)
	if err != nil {
		switch err.Error() {
		case "1":
			c.JSON(http.StatusNotFound, gin.H{"error": "Username cannot empty or at lease 6 character"})
		case "2":
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not valid email"})
		case "3":
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not empty password and password must have at least 10 characters"})
		case "4":
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Already Exit,choose another one"})
		case "5":
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Can not hash password"})
		case "0":
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Repassword is not same with password"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  res.Message,
		"userid":   res.UserID,
		"username": res.Username,
		"email":    res.Email,
	})
}
