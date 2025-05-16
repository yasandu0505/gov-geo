package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

var AppConfig Config

func Load() {
	// Load .env file if present
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables.")
	}

	AppConfig = Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "yourpassword"),
		DBName:     getEnv("DB_NAME", "yourdb"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key string, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
