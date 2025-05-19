package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"go-mysql-backend/config"
	"go-mysql-backend/internal/handlers"
	"go-mysql-backend/internal/repository"
	"go-mysql-backend/internal/service"
	"go-mysql-backend/routes"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL not provided")
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Successfully connected to PostgreSQL!")

	orgRepo := repository.NewOrganizationRepository(db)
	orgService := service.NewOrganizationService(orgRepo)
	orgHandler := handlers.NewOrganizationHandler(orgService)

	router := mux.NewRouter()
	routes.SetupOrgRoutes(router, orgHandler)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
