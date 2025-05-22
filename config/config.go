package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	UseNeo4j    bool
	Neo4jURI    string
	Neo4jUser   string
	Neo4jPass   string
}

func LoadConfig() Config {
	// Try to load .env file from the project root
	workingDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting working directory: %v", err)
	} else {
		log.Printf("Working directory: %s", workingDir)
	}

	// Try different locations for the .env file
	envLocations := []string{
		".env",
		filepath.Join(workingDir, ".env"),
		filepath.Join(workingDir, "..", ".env"),
	}

	var loaded bool
	for _, location := range envLocations {
		if err := godotenv.Load(location); err == nil {
			log.Printf("Loaded configuration from %s", location)
			loaded = true
			break
		}
	}

	if !loaded {
		log.Println("No .env file found, using environment variables")
	}

	// Debug: Print environment variables
	fmt.Println("Environment variables:")
	fmt.Printf("USE_NEO4J=%s\n", os.Getenv("USE_NEO4J"))
	fmt.Printf("NEO4J_URI=%s\n", os.Getenv("NEO4J_URI"))
	fmt.Printf("NEO4J_USERNAME=%s\n", os.Getenv("NEO4J_USERNAME"))
	fmt.Printf("DATABASE_URL=%s\n", os.Getenv("DATABASE_URL"))

	useNeo4j := os.Getenv("USE_NEO4J") == "true"

	// Debug: Print config values
	config := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		UseNeo4j:    useNeo4j,
		Neo4jURI:    getEnvWithDefault("NEO4J_URI", "neo4j://localhost:7687"),
		Neo4jUser:   getEnvWithDefault("NEO4J_USERNAME", "neo4j"),
		Neo4jPass:   os.Getenv("NEO4J_PASSWORD"),
	}

	fmt.Printf("Config loaded: UseNeo4j=%v, Neo4jURI=%s, Neo4jUser=%s\n",
		config.UseNeo4j, config.Neo4jURI, config.Neo4jUser)

	return config
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
