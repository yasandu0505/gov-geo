package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Neo4jInterface provides a wrapper around Neo4j operations
type Neo4jInterface struct {
	driver neo4j.Driver
}

// NewNeo4jInterface creates a new Neo4j interface
func NewNeo4jInterface(driver neo4j.Driver) *Neo4jInterface {
	return &Neo4jInterface{
		driver: driver,
	}
}

// Close closes the Neo4j driver
func (n *Neo4jInterface) Close() error {
	return n.driver.Close()
}

// ExecuteQuery executes a Cypher query with parameters and returns the result
func (n *Neo4jInterface) ExecuteQuery(query string, params map[string]interface{}) ([]map[string]interface{}, error) {
	// Create a new session
	session := n.driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close()

	// Execute the query
	result, err := session.Run(query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to run query: %w", err)
	}

	// Collect all records
	records, err := result.Collect()
	if err != nil {
		return nil, fmt.Errorf("failed to collect results: %w", err)
	}

	// Convert records to maps
	var resultMaps []map[string]interface{}
	for _, record := range records {
		recordMap := make(map[string]interface{})
		for i, key := range record.Keys {
			recordMap[key] = record.Values[i]
		}
		resultMaps = append(resultMaps, recordMap)
	}

	return resultMaps, nil
}

// Ministry represents a government ministry
type Ministry struct {
	ID                int
	Name              string
	Google_map_script string
}

// Department represents a government department
type Department struct {
	ID                int
	Name              string
	Google_map_script string
	MinistryID        int
}

// createMinistry creates a ministry node and links it to the government
func createMinistry(neo4jInterface *Neo4jInterface, ministry Ministry) ([]map[string]interface{}, error) {
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

func main() {
	fmt.Println("Adding new ministry to Neo4j database...")

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: No .env file found. Using system environment variables.")
	} else {
		fmt.Println("Loaded environment variables from .env file")
	}

	// Get Neo4j connection details from environment variables
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
	neo4jInterface := NewNeo4jInterface(driver)
	defer func() {
		if err := neo4jInterface.Close(); err != nil {
			log.Printf("Error closing Neo4j connection: %v", err)
		} else {
			fmt.Println("Neo4j connection closed successfully.")
		}
	}()

	// Define a new ministry - using ID 5 since we already have ministries with IDs 1-4
	ministry := Ministry{
		ID:                5,
		Name:              "Ministry of Agriculture",
		Google_map_script: "<script>agriculture map</script>",
	}

	// Create the ministry
	startTime := time.Now()
	fmt.Printf("Creating ministry: %s...\n", ministry.Name)
	result, err := createMinistry(neo4jInterface, ministry)
	if err != nil {
		log.Fatalf("Error creating ministry: %v", err)
	}
	fmt.Printf("Ministry '%s' created successfully\n", ministry.Name)

	// Print ministry details
	for _, record := range result {
		fmt.Printf("Node created: %v\n", record["m"])
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\nTotal execution time: %v\n", elapsed)
}
