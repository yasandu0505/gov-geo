package db

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func InitNeo4j() (neo4j.DriverWithContext, error) {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found. Using system environment variables.")
	}

	dbUri := os.Getenv("NEO4J_URL")
	dbUser := os.Getenv("NEO4J_USER")
	dbPassword := os.Getenv("NEO4J_PASSWORD")

	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		return nil, err
	}

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		driver.Close(ctx)
		return nil, err
	}

	log.Println("✅ Connected to Neo4j")
	return driver, nil
}
