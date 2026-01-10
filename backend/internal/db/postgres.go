package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
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

func ConnectPostgresPool() *pgxpool.Pool {
	cfg := config.LoadConfig()
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("❌ Failed to create pgxpool: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}

	log.Println("✅ Postgres pool connection established")
	return pool
}
