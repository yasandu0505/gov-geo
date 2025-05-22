package db

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func InitNeo4j() {
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
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("✅ Connected to Neo4j")
}
