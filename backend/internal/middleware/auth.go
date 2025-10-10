package middleware

import (
	"net/http"

	authService "app/internal/app/auth/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Service authService.AuthService
}

func (handler *AuthHandler) AuthenticationMiddleware(c *gin.Context) {
	tokenString, err := c.Cookie("token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid authentication cookie"})
		return
	}

	claims, err := handler.Service.ValidateJWT(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	c.Set("userClaims", claims)
	c.Next()
}
