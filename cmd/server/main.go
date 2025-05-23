package main

import (
	"fmt"
	"log"
	"net/http"

	"go-mysql-backend/config"
	"go-mysql-backend/internal/db"
	"go-mysql-backend/internal/handlers"
	"go-mysql-backend/internal/repository"
	"go-mysql-backend/internal/service"
	"go-mysql-backend/routes"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {

	dbType := config.LoadType()
	if dbType == "postgres" {

		db := db.InitPostgres()

		orgRepo := repository.NewOrganizationRepository(db)
		orgService := service.NewOrganizationService(orgRepo)
		orgHandler := handlers.NewOrganizationHandler(orgService)

		router := mux.NewRouter()
		routes.SetupOrgRoutes(router, orgHandler)

		startServer(router)

	} else if dbType == "neo4j" {

		neo4jDriver, err := db.InitNeo4j()
		if err != nil {
			log.Fatal("Failed to connect to Neo4j:", err)
		}
		neoRepo := repository.NewNeo4jRepository(neo4jDriver)
		neoService := service.NewNeo4JService(neoRepo)
		neoHandler := handlers.NewNeo4JHandler(neoService)
		router := mux.NewRouter()
		routes.SetupNeo4JRoutes(router, neoHandler)

		startServer(router)
	}
}

func startServer(router *mux.Router) {
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler(router)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
