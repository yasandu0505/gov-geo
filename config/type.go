package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadType() string {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return os.Getenv("DATABASE_TYPE")

}
