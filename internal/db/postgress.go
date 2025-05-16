package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitPostgres() *sql.DB {
	// Load .env (optional if already loaded in main)
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found. Using system environment variables.")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL not set in environment variables")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to open DB connection: %v", err)
	}

	// Verify DB connection
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Cannot connect to DB: %v", err)
	}

	fmt.Println("✅ Connected to PostgreSQL")
	return db
}
