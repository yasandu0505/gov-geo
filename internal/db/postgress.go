package db

import (
	"database/sql"
	"fmt"
	"log"

	"go-mysql-backend/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	cfg := config.AppConfig

	var err error
	DB, err = sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("❌ Database unreachable:", err)
	}

	fmt.Println("✅ PostgreSQL connected successfully!")
}
