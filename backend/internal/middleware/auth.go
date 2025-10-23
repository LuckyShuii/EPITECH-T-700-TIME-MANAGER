package middleware

import (
	"net/http"
	"slices"

	authService "app/internal/app/auth/service"

	Config "app/internal/config"

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

func (handler *AuthHandler) RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve claims from context
		claims, exists := c.Get("userClaims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": Config.ErrorMessages()["NO_CLAIMS"]})
			return
		}

		// check the roles in the claims
		userRoles := claims.(*authService.Claims).Roles

		// Check if user has at least one of the required roles
		for _, requiredRole := range roles {
			if slices.Contains(userRoles, requiredRole) || slices.Contains(userRoles, "admin") || requiredRole == "all" {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
	}
}
