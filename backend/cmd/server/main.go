package main

import (
	"app/internal/app/mailer"
	"app/internal/app/mailer/provider"
	"app/internal/config"
	"app/internal/db"
	"app/internal/router"
	"context"
	"log"
	"time"

	"app/cmd/server/docs"

	"github.com/gin-contrib/cors"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	mailerservice "app/internal/app/mailer/service"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig()

	mailerProvider := &provider.BrevoMailer{
		APIKey: cfg.Mail.APIKey,
	}

	mailer.Service = mailerservice.NewMailerService(mailerProvider)

	// Seed database with fixtures in DEV mode
	ctx := context.Background()
	pool := db.ConnectPostgresPool()
	defer pool.Close()

	fixturesPath := "/app/fixtures.sql"
	if err := db.SeedIfEmptyUsersDevOnly(ctx, pool, fixturesPath, 1, cfg.ProjectStatus); err != nil {
		log.Printf("⚠️  Warning: Failed to seed database: %v", err)
	} else {
		log.Println("✅ Database seeding completed successfully")
	}

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

	// ✅ Configuration Swagger
	docs.SwaggerInfo.Title = "Time Manager API"
	docs.SwaggerInfo.Description = "API documentation for Time Manager backend"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = "localhost:8081"

	// ✅ Route Swagger (accessible via /api/swagger/index.html)
	r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":5000")
}
