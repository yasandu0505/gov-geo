package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
}

var AppConfig Config

func Load() {
	// Load .env file if present
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables.")
	}

	AppConfig = Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
	}

	if AppConfig.DatabaseURL == "" {
		log.Fatal("❌ DATABASE_URL is required but not set.")
	}
}

func getEnv(key string, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
