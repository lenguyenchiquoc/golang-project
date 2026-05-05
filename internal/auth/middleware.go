package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ErrMissingAuth   = "AUTH_MISSING_HEADER"
	ErrInvalidFormat = "AUTH_INVALID_FORMAT"
	ErrTokenExpired  = "AUTH_TOKEN_EXPIRED"
	ErrTokenInvalid  = "AUTH_TOKEN_INVALID"
)

func JWTMiddleware(service *AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrMissingAuth})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidFormat})
			return
		}

		tokenStr := parts[1]
		claims, err := service.ValidateJWT(tokenStr)
		if err != nil {
			errMsg := ErrTokenInvalid
			if strings.Contains(err.Error(), "expired") {
				errMsg = ErrTokenExpired
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   errMsg,
				"details": err.Error(),
			})
			return
		}
		c.Set("user_id", claims["user_id"])
		c.Set("username", claims["username"])

		c.Next()
	}
}
