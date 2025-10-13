package main

import (
	"app/internal/config"
	"app/internal/router"
	"time"

	"github.com/gin-contrib/cors"
)

func main() {
	cfg := config.LoadConfig()

	r := router.SetupRouter()

	allowedOrigins := []string{
		"http://localhost:5173",
		cfg.FrontendURL,
		"http://frontend:5173",
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Run(":5000")
}
