package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func ConnectRedis() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	if host == "" {
		host = "redis"
	}
	if port == "" {
		port = "6379"
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,
	})

	// Connection test
	_, err := client.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("❌ Connection to Redis failed: %v", err)
	}

	log.Println("✅ Connection to Redis successful:", addr)
	RedisClient = client
	return client
}
