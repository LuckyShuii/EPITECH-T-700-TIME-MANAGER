package main

import (
	"net/http"
	"app/src/ping"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/api/ping", ping.Ping)

	router.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API Working :)",
		})
	})
	
	router.Run(":5000")
}
