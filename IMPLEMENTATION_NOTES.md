# Neo4j Integration Implementation Notes

## Overview

This project has been extended to support both PostgreSQL and Neo4j as data stores. The implementation follows a clean architecture pattern with proper separation of concerns.

## Key Components

1. **Database Connection**:
   - `internal/db/neo4j.go`: Handles Neo4j driver initialization
   - `internal/db/neo4j_interface.go`: Provides a utility wrapper around Neo4j operations

2. **Repository Layer**:
   - `internal/repository/organization_interface.go`: Defines a common interface for all organization repositories
   - `internal/repository/organization.go`: PostgreSQL implementation
   - `internal/repository/organization_neo4j.go`: Neo4j implementation

3. **Data Loading Utility**:
   - `cmd/neo4j/main.go`: Go utility to initialize Neo4j with basic data structure

4. **Python Integration**:
   - `neo4j_util/neo4j_interface.py`: Python Neo4j interface class
   - `neo4j_utils_main.py`: Python script to add government node data to Neo4j

## Configuration

The system uses environment variables to determine which database to use:

```
USE_NEO4J=true  # Set to true to use Neo4j, false for PostgreSQL
NEO4J_URI=neo4j+s://08141f67.databases.neo4j.io
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=4_P7gGPQWUO9n8bv8-lZlhSU2ofnjNXl1EKJyrh8Zx8
DATABASE_URL=postgres://user:password@localhost:5432/orgdb  # Only used if USE_NEO4J=false
```

## Data Model

The data model represents a government organization hierarchy:

- Government: The top-level entity
- Ministries: Entities that belong to the government
- Departments: Entities that belong to ministries

In Neo4j, these are represented as nodes with relationships:
- (government)-[:HAS_MINISTRY]->(ministry)
- (ministry)-[:HAS_DEPARTMENT]->(department)

## Endpoints

The API provides these endpoints regardless of the database backend used:

- GET /ministries - List all ministries with their departments
- GET /departments - List all departments
- POST /ministries - Create a new ministry
- POST /departments - Create a new department
- GET /ministries/{id} - Get a specific ministry
- GET /departments/{id} - Get a specific department

## Running the Application

1. Set up your environment variables in a `.env` file
2. Run the server: `go run cmd/server/main.go`
3. If using Neo4j, initialize data: `go run cmd/neo4j/main.go`
4. Alternative Python Neo4j initialization: `python neo4j_utils_main.py` 