# Script to temporarily set environment variables for PostgreSQL testing
# Keep USE_NEO4J=true but add DATABASE_URL so you can easily switch between databases

# Set the DATABASE_URL environment variable
$env:DATABASE_URL = "postgres://user:password@localhost:5432/orgdb"

# Display current environment settings
Write-Host "Current environment settings:"
Write-Host "USE_NEO4J: $env:USE_NEO4J"
Write-Host "DATABASE_URL: $env:DATABASE_URL"
Write-Host "NEO4J_URI: $env:NEO4J_URI"
Write-Host "NEO4J_USERNAME: $env:NEO4J_USERNAME"

Write-Host "`nTo switch to PostgreSQL, use:"
Write-Host '$env:USE_NEO4J = "false"'
Write-Host "go run cmd/server/main.go"

Write-Host "`nTo switch back to Neo4j, use:"
Write-Host '$env:USE_NEO4J = "true"'
Write-Host "go run cmd/server/main.go" 