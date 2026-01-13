package config

import (
	"os"
)

var (
	Port         = getEnv("PORT", "8080")
	DatabasePath = getEnv("DATABASE_PATH", "./fired_calendar.db")
	SessionKey   = getEnv("SESSION_KEY", "your-secret-key-change-this-in-production")
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
