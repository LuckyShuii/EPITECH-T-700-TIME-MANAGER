/**
 * Configuration package for managing application settings such as environment variables, ...
 */
package config

import (
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
	ProjectStatus      string
	FrontendURL        string
	RedisHost          string
	RedisPort          string
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
		ProjectStatus:      getEnv("PROJECT_STATUS", os.Getenv("PROJECT_STATUS")),
		FrontendURL:        getEnv("FRONTEND_URL", os.Getenv("FRONTEND_URL")),
		RedisHost:          getEnv("REDIS_HOST", os.Getenv("REDIS_HOST")),
		RedisPort:          getEnv("REDIS_PORT", os.Getenv("REDIS_PORT")),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func ErrorMessages() map[string]string {
	return map[string]string{
		"NO_CLAIMS":             "missing claims",
		"INVALID_REQUEST":       "invalid request",
		"WEEKLY_RATE_NOT_FOUND": "failed to find weekly rate",
	}
}
