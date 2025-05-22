package main

import (
	"fmt"
	"log"

	"go-mysql-backend/neo4j_util_go"
)

func main() {
	fmt.Println("Neo4j Utility - Government Node Creator")
	fmt.Println("---------------------------------------")

	// Run the Neo4j utility
	neo4j_util_go.RunNeo4jUtil()

	log.Println("Neo4j utility completed")
}
