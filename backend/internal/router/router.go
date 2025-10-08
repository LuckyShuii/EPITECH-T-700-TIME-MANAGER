/**
 * Router for handling API requests and defining route handlers
 * For now all in one file, then start the structure in this folder once basic concepts are understood
 */

package router

import (
	"net/http"
	"app/internal/app/user"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/api/ping", user.Ping)
	router.GET("/api/users", user.ReturnFakeUsers)
	router.GET("/api/user/:id", user.ReturnFakeUserById)
	router.POST("/api/users", user.AddUser)

	router.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API Working :)"})
	})

	router.GET("/api/example", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Example route working! :)"})
	})

	return router
}