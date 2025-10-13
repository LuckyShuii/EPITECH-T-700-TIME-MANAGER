package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"app/internal/config"
)

func ConnectPostgres() *gorm.DB {
	host := config.LoadConfig().DBHost
	user := config.LoadConfig().DBUser
	password := config.LoadConfig().DBPassword
	dbname := config.LoadConfig().DBName
	port := config.LoadConfig().DBPort

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Paris",
		host, user, password, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("❌ Failed to connect to Postgres: %v", err)
	}

	log.Println("✅ Connected to Postgres successfully")
	return db
}
