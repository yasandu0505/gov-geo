package neo4j_util_go

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Neo4jInterface provides an interface to interact with Neo4j database
type Neo4jInterface struct {
	driver neo4j.Driver
}

// NewNeo4jInterface initializes a new Neo4j connection
func NewNeo4jInterface() (*Neo4jInterface, error) {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found. Using system environment variables.")
	}

	// Get Neo4j connection details
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
		log.Println("Warning: NEO4J_PASSWORD not set. Using empty password.")
	}

	// Create Neo4j driver
	driver, err := neo4j.NewDriver(
		uri,
		neo4j.BasicAuth(username, password, ""),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Neo4j driver: %w", err)
	}

	// Verify connection
	err = driver.VerifyConnectivity()
	if err != nil {
		// Close the driver if verification fails
		driver.Close()
		return nil, fmt.Errorf("failed to connect to Neo4j: %w", err)
	}

	fmt.Println("✅ Connected to Neo4j")

	return &Neo4jInterface{
		driver: driver,
	}, nil
}

// Close closes the Neo4j connection
func (n *Neo4jInterface) Close() {
	if n.driver != nil {
		n.driver.Close()
		fmt.Println("Neo4j connection closed.")
	}
}

// ExecuteQuery executes a Cypher query with parameters and returns the result
func (n *Neo4jInterface) ExecuteQuery(query string, params map[string]interface{}) ([]map[string]interface{}, error) {
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

// ExecuteQueryWithRetry executes a Cypher query with retry mechanism
func (n *Neo4jInterface) ExecuteQueryWithRetry(query string, params map[string]interface{}, maxAttempts int) ([]map[string]interface{}, error) {
	if params == nil {
		params = make(map[string]interface{})
	}

	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("Attempt %d of %d to execute query...\n", attempt, maxAttempts)

		result, err := n.ExecuteQuery(query, params)
		if err == nil {
			fmt.Println("Query executed successfully.")
			return result, nil
		}

		lastErr = err
		fmt.Printf("Error on attempt %d: %v\n", attempt, err)

		if attempt < maxAttempts {
			// Exponential backoff
			waitTime := time.Duration(2*attempt) * time.Second
			fmt.Printf("Waiting %v seconds before retrying...\n", waitTime.Seconds())
			time.Sleep(waitTime)
		}
	}

	return nil, fmt.Errorf("maximum retry attempts reached. Last error: %w", lastErr)
}
