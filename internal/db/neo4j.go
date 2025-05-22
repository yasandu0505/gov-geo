package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// InitNeo4j initializes and returns a Neo4j driver
func InitNeo4j() (neo4j.Driver, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found. Using system environment variables.")
	}

	// Get Neo4j connection details from environment variables
	uri := os.Getenv("NEO4J_URI")
	if uri == "" {
		uri = "neo4j://localhost:7687" // Default URI
	}

	username := os.Getenv("NEO4J_USERNAME")
	if username == "" {
		username = "neo4j" // Default username
	}

	password := os.Getenv("NEO4J_PASSWORD")
	if password == "" {
		log.Println("⚠️ NEO4J_PASSWORD not set. Using empty password.")
	}

	// Create Neo4j driver
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, fmt.Errorf("failed to create Neo4j driver: %w", err)
	}

	// Verify connection
	err = driver.VerifyConnectivity()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Neo4j: %w", err)
	}

	fmt.Println("✅ Connected to Neo4j")
	return driver, nil
}
