package db

import (
	"fmt"
	"time"

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

// ExecuteQueryWithRetry executes a Cypher query with retry mechanism
func (n *Neo4jInterface) ExecuteQueryWithRetry(query string, params map[string]interface{}, maxAttempts int) ([]map[string]interface{}, error) {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("Attempt %d of %d to execute query...\n", attempt, maxAttempts)

		result, err := n.ExecuteQuery(query, params)
		if err == nil {
			fmt.Println("Query executed successfully.")
			return result, nil
		}

		lastErr = err
		fmt.Printf("Error on attempt %d: %s\n", attempt, err)

		if attempt < maxAttempts {
			// Exponential backoff
			waitTime := time.Duration(2*attempt) * time.Second
			fmt.Printf("Waiting %v seconds before retrying...\n", waitTime.Seconds())
			time.Sleep(waitTime)
		}
	}

	return nil, fmt.Errorf("maximum retry attempts reached. Last error: %w", lastErr)
}
