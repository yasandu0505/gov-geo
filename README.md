# Government Geographic Data API

A robust RESTful API built in Go (Golang) that provides access to government geographic data, supporting both PostgreSQL and Neo4j databases. The project follows clean architecture principles with proper separation of concerns.

##  Project Structure

```
gov-geo/
│
├── cmd/
│   └── server/             # Main entry point of the application
│
├── config/                 # Configuration settings
│
├── internal/
│   ├── db/                 # Database connection logic
│   ├── models/             # Data models (structs)
│   ├── repository/         # Data access layer
│   ├── service/            # Business logic layer
│   └── handler/            # HTTP request handlers
│
├── routes/                 # Route definitions
│
├── go.mod                  # Go module file
├── go.sum                  # Go module checksum file
└── README.md              # Project documentation
```

## Technologies Used

- Go (Golang) 1.24.3
- PostgreSQL
- Neo4j
- Gorilla Mux (HTTP router)
- CORS support
- Environment variable management (godotenv)

##  Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/yasandu0505/gov-geo.git
cd gov-geo
```

### 2. Set up the databases

#### PostgreSQL Setup
```sql
CREATE DATABASE gov_geo;
```

#### Neo4j Setup
- Install Neo4j Desktop or use Neo4j Aura
- Create a new database
- Note down your connection credentials

### 3. Configure environment variables

Create a `.env` file in the root directory with the following variables:

```env
# PostgreSQL Configuration
POSTGRES_USER=your_username
POSTGRES_PASSWORD=your_password
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=gov_geo

# Neo4j Configuration
NEO4J_URI=bolt://localhost:7687
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=your_password
```

### 4. Install dependencies

```bash
go mod download
```

### 5. Run the server

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

##  API Endpoints

The API provides various endpoints for accessing and manipulating geographic data. Detailed API documentation can be found in the routes directory.

##  Security

- CORS is enabled for secure cross-origin requests
- Environment variables are used for sensitive configuration
- Database credentials are never hardcoded

## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

