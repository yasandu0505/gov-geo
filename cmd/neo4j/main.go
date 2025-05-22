package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"go-mysql-backend/internal/db"
	"go-mysql-backend/internal/models"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// createGovernmentNode creates a government node with retry mechanism
func createGovernmentNode(neo4jInterface *db.Neo4jInterface) ([]map[string]interface{}, error) {
	query := `
	MERGE (g:government {id: $id})
	SET g.name = $name
	RETURN g
	`
	params := map[string]interface{}{
		"id":   "gov_01",
		"name": "Government of Sri Lanka",
	}

	maxAttempts := 3
	return neo4jInterface.ExecuteQueryWithRetry(query, params, maxAttempts)
}

// createMinistry creates a ministry node and links it to the government
func createMinistry(neo4jInterface *db.Neo4jInterface, ministry models.Ministry) ([]map[string]interface{}, error) {
	query := `
	MATCH (g:government {id: $govID})
	MERGE (m:ministry {id: $id, name: $name})
	SET m.google_map_script = $mapScript
	MERGE (g)-[:HAS_MINISTRY]->(m)
	RETURN m
	`
	params := map[string]interface{}{
		"govID":     "gov_01",
		"id":        fmt.Sprintf("ministry_%d", ministry.ID),
		"name":      ministry.Name,
		"mapScript": ministry.Google_map_script,
	}

	return neo4jInterface.ExecuteQuery(query, params)
}

// createDepartment creates a department node and links it to its ministry
func createDepartment(neo4jInterface *db.Neo4jInterface, dept models.Department) ([]map[string]interface{}, error) {
	query := `
	MATCH (m:ministry {id: $ministryID})
	MERGE (d:department {id: $id, name: $name})
	SET d.google_map_script = $mapScript
	MERGE (m)-[:HAS_DEPARTMENT]->(d)
	RETURN d
	`
	params := map[string]interface{}{
		"ministryID": fmt.Sprintf("ministry_%d", dept.MinistryID),
		"id":         fmt.Sprintf("dept_%d", dept.ID),
		"name":       dept.Name,
		"mapScript":  dept.Google_map_script,
	}

	return neo4jInterface.ExecuteQuery(query, params)
}

func main() {
	// Try loading the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: No .env file found. Using system environment variables.")
	} else {
		fmt.Println("Loaded environment variables from .env file")
	}

	// Print environment variables for debugging
	fmt.Println("Environment variables:")
	fmt.Printf("NEO4J_URI=%s\n", os.Getenv("NEO4J_URI"))
	fmt.Printf("NEO4J_USERNAME=%s\n", os.Getenv("NEO4J_USERNAME"))

	// Get connection details from environment variables
	uri := os.Getenv("NEO4J_URI")
	if uri == "" {
		uri = "neo4j://localhost:7687" // Default URI
		fmt.Println("Warning: NEO4J_URI not set, using default: " + uri)
	}

	username := os.Getenv("NEO4J_USERNAME")
	if username == "" {
		username = "neo4j" // Default username
		fmt.Println("Warning: NEO4J_USERNAME not set, using default: " + username)
	}

	password := os.Getenv("NEO4J_PASSWORD")
	if password == "" {
		fmt.Println("Error: NEO4J_PASSWORD not set")
		log.Fatal("NEO4J_PASSWORD environment variable is required")
	}

	// Initialize Neo4j driver
	fmt.Printf("Connecting to Neo4j: %s with user %s\n", uri, username)
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
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
	defer func() {
		if err := neo4jInterface.Close(); err != nil {
			log.Printf("Error closing Neo4j connection: %v", err)
		} else {
			fmt.Println("Neo4j connection closed successfully.")
		}
	}()

	// Create the government node
	startTime := time.Now()
	fmt.Println("Creating government node...")
	govResult, err := createGovernmentNode(neo4jInterface)
	if err != nil {
		log.Fatalf("Error creating government node: %v", err)
	}

	fmt.Println("Government node created successfully:")
	for _, record := range govResult {
		fmt.Printf("Node: %v\n", record["g"])
	}

	// Create sample ministry
	ministry := models.Ministry{
		ID:                1,
		Name:              "Ministry of Digital Infrastructure",
		Google_map_script: "<script>map code here</script>",
	}

	fmt.Println("\nCreating ministry node...")
	ministryResult, err := createMinistry(neo4jInterface, ministry)
	if err != nil {
		log.Fatalf("Error creating ministry node: %v", err)
	}

	fmt.Println("Ministry node created successfully:")
	for _, record := range ministryResult {
		fmt.Printf("Node: %v\n", record["m"])
	}

	// Create sample department
	department := models.Department{
		ID:                1,
		Name:              "Department of Computer Services",
		MinistryID:        1,
		Google_map_script: "<script>map code here</script>",
	}

	fmt.Println("\nCreating department node...")
	deptResult, err := createDepartment(neo4jInterface, department)
	if err != nil {
		log.Fatalf("Error creating department node: %v", err)
	}

	fmt.Println("Department node created successfully:")
	for _, record := range deptResult {
		fmt.Printf("Node: %v\n", record["d"])
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\nTotal execution time: %v\n", elapsed)
}
