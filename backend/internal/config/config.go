/**
 * Configuration package for managing application settings such as environment variables, ...
 */
package config

import (
	"log"
	"os"
)

type Config struct {
	Port               string
	DBHost             string
	DBUser             string
	DBPassword         string
	DBName             string
	DBPort             string
	JWTSecret          string
	JWTExpirationHours string
}

func LoadConfig() *Config {
	config := &Config{
		Port:               getEnv("PORT", os.Getenv("DB_PORT")),
		DBHost:             getEnv("DB_HOST", os.Getenv("DB_HOST")),
		DBUser:             getEnv("DB_USER", os.Getenv("DB_USER")),
		DBPassword:         getEnv("DB_PASSWORD", os.Getenv("DB_PASSWORD")),
		DBName:             getEnv("DB_DATABASE", os.Getenv("DB_DATABASE")),
		DBPort:             getEnv("DB_PORT", os.Getenv("DB_PORT")),
		JWTSecret:          getEnv("JWT_SECRET", os.Getenv("JWT_SECRET")),
		JWTExpirationHours: getEnv("JWT_EXPIRATION_HOURS", os.Getenv("JWT_EXPIRATION_HOURS")),
	}

	log.Printf("Configuration loaded: %+v\n", config)
	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
