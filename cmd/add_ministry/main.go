package main

import (
	"fmt"
	"log"
	"time"

	"go-mysql-backend/neo4j_util_go"
)

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
func createMinistry(neo4jInterface *neo4j_util_go.Neo4jInterface, ministry Ministry) ([]map[string]interface{}, error) {
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

	// Initialize Neo4j interface using our new package
	neo4jInterface, err := neo4j_util_go.NewNeo4jInterface()
	if err != nil {
		log.Fatalf("Failed to initialize Neo4j interface: %v", err)
	}
	defer neo4jInterface.Close()

	// Define a new ministry - using ID 6 since we already created one with ID 5
	ministry := Ministry{
		ID:                6,
		Name:              "Ministry of Digital Infrastructure",
		Google_map_script: "<script>digital infrastructure map</script>",
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
