package neo4j_util_go

import (
	"fmt"
	"log"
	"time"
)

// CreateGovernmentNode creates a government node with retry mechanism
func CreateGovernmentNode(neo4j *Neo4jInterface) ([]map[string]interface{}, error) {
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
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		fmt.Printf("Attempt %d of %d to create government node...\n", attempt, maxAttempts)
		result, err := neo4j.ExecuteQuery(query, params)
		if err == nil {
			fmt.Println("Successfully created government node.")
			return result, nil
		}

		lastErr = err
		fmt.Printf("Error on attempt %d: %v\n", attempt, err)

		if attempt < maxAttempts {
			waitTime := time.Duration(2*attempt) * time.Second
			fmt.Printf("Waiting %v seconds before retrying...\n", waitTime.Seconds())
			time.Sleep(waitTime)
		}
	}

	return nil, fmt.Errorf("maximum retry attempts reached. Failed to create government node: %w", lastErr)
}

// RunNeo4jUtil is the main entry point for the Neo4j utility
func RunNeo4jUtil() {
	// Initialize Neo4j interface
	neo4jInterface, err := NewNeo4jInterface()
	if err != nil {
		log.Fatalf("Failed to initialize Neo4j interface: %v", err)
	}

	// Ensure we always close the connection
	defer neo4jInterface.Close()

	// Create the government node
	govNode, err := CreateGovernmentNode(neo4jInterface)
	if err != nil {
		log.Printf("Error creating government node: %v", err)
		return
	}

	// Print government node details
	fmt.Println("Government node details:")
	for _, record := range govNode {
		fmt.Printf("Node: %v\n", record["g"])
	}
}
