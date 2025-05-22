package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"go-mysql-backend/internal/db"
	"go-mysql-backend/internal/handlers"
	"go-mysql-backend/internal/repository"
	"go-mysql-backend/internal/service"
	"go-mysql-backend/routes"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	// Debug: Print environment values
	fmt.Println("Starting server...")

	// Try loading the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: No .env file found. Using system environment variables.")
	} else {
		fmt.Println("Loaded environment variables from .env file")
	}

	// Print environment variables for debugging
	fmt.Println("Environment variables:")
	fmt.Printf("USE_NEO4J=%s\n", os.Getenv("USE_NEO4J"))
	fmt.Printf("NEO4J_URI=%s\n", os.Getenv("NEO4J_URI"))
	fmt.Printf("NEO4J_USERNAME=%s\n", os.Getenv("NEO4J_USERNAME"))
	fmt.Printf("DATABASE_URL=%s\n", os.Getenv("DATABASE_URL"))

	// Determine which database to use
	useNeo4j := os.Getenv("USE_NEO4J") == "true"
	fmt.Printf("Using Neo4j: %v\n", useNeo4j)

	// Get database connection details
	neo4jURI := os.Getenv("NEO4J_URI")
	if neo4jURI == "" {
		neo4jURI = "neo4j://localhost:7687" // Default URI
	}

	neo4jUsername := os.Getenv("NEO4J_USERNAME")
	if neo4jUsername == "" {
		neo4jUsername = "neo4j" // Default username
	}

	neo4jPassword := os.Getenv("NEO4J_PASSWORD")
	databaseURL := os.Getenv("DATABASE_URL")

	// Check if we have the required configuration
	if databaseURL == "" && !useNeo4j {
		log.Fatal("DATABASE_URL not provided and USE_NEO4J is not set to true")
	}

	if useNeo4j && neo4jPassword == "" {
		log.Fatal("NEO4J_PASSWORD not provided but USE_NEO4J is set to true")
	}

	var orgRepo repository.OrganizationRepositoryInterface

	if useNeo4j {
		// Use Neo4j as the data store
		fmt.Printf("Attempting to connect to Neo4j: %s with user %s\n", neo4jURI, neo4jUsername)
		driver, err := neo4j.NewDriver(
			neo4jURI,
			neo4j.BasicAuth(neo4jUsername, neo4jPassword, ""),
		)
		if err != nil {
			log.Fatalf("Failed to create Neo4j driver: %v", err)
		}
		defer driver.Close()

		// Verify connection
		fmt.Println("Verifying Neo4j connectivity...")
		err = driver.VerifyConnectivity()
		if err != nil {
			log.Fatalf("Failed to connect to Neo4j: %v", err)
		}

		fmt.Println("✅ Successfully connected to Neo4j!")

		// Create Neo4j interface
		neo4jInterface := db.NewNeo4jInterface(driver)
		orgRepo = repository.NewOrganizationNeo4jRepository(neo4jInterface)
	} else {
		// Use PostgreSQL as the data store
		fmt.Println("Attempting to connect to PostgreSQL...")
		sqlDB, err := sql.Open("postgres", databaseURL)
		if err != nil {
			log.Fatalf("Failed to open PostgreSQL connection: %v", err)
		}
		defer sqlDB.Close()

		err = sqlDB.Ping()
		if err != nil {
			log.Fatalf("Failed to ping PostgreSQL: %v", err)
		}

		fmt.Println("✅ Successfully connected to PostgreSQL!")
		orgRepo = repository.NewOrganizationRepository(sqlDB)
	}

	orgService := service.NewOrganizationService(orgRepo)
	orgHandler := handlers.NewOrganizationHandler(orgService)

	router := mux.NewRouter()
	routes.SetupOrgRoutes(router, orgHandler)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
