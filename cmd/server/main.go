package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get database URI from env
	dbURI := os.Getenv("DATABASE_URL")
	if dbURI == "" {
		log.Fatal("DATABASE_URL not found in .env file")
	}

	// Connect to PostgreSQL
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Ping the database to verify connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Database unreachable:", err)
	}

	fmt.Println("✅ Successfully connected to PostgreSQL!")
	defer db.Close()

	// You can now use `db` to perform queries
}
