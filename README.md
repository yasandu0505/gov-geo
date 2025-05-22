# Go Backend API

A modular RESTful API built in Go (Golang) that connects to either a PostgreSQL database or Neo4j graph database. It follows a clean architecture with proper separation of concerns (handlers, services, repositories, and models).

## 📁 Project Structure

go-backend/
│
├── cmd/
│   ├── server/             # Main entry point of the application
│   └── neo4j/              # Neo4j data loading utility
│
├── config/                 # Configuration settings (e.g. database credentials)
│
├── internal/
│   ├── db/                 # Database connection logic
│   │   ├── postgress.go    # PostgreSQL connection
│   │   ├── neo4j.go        # Neo4j connection
│   │   └── neo4j_interface.go # Neo4j query interface
│   ├── models/             # Data models (structs)
│   ├── repository/         # Data access layer (SQL and Neo4j queries)
│   ├── service/            # Business logic layer
│   └── handler/            # HTTP request handlers (controllers)
│
├── routes/                 # Route definitions
│
├── go.mod                  # Go module file
├── go.sum                  # Go module checksum file
└── README.md               # Project documentation

## 🛠️ Technologies Used

- Go (Golang)
- PostgreSQL or Neo4j
- Gorilla Mux (HTTP router)
- Git

## 🚀 Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/yasandu0505/gov-geo.git
cd gov-geo
```

### 2. Set up the database

#### Option 1: PostgreSQL

```sql
CREATE DATABASE orgdb;
USE orgdb;

CREATE TABLE ministry (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  google_map_script TEXT
);

CREATE TABLE department (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  ministry_id INT REFERENCES ministry(id),
  google_map_script TEXT
);
```

#### Option 2: Neo4j

1. Install Neo4j Desktop or use Neo4j Aura cloud service
2. Create a new database
3. Set the username and password
4. The Neo4j schema will be created automatically when running the application

### 3. Configure environment variables

Create a `.env` file in the project root:

```
# Database configuration
# PostgreSQL connection string
DATABASE_URL=postgres://user:password@localhost:5432/orgdb

# Neo4j configuration
USE_NEO4J=true  # Set to true to use Neo4j instead of PostgreSQL
NEO4J_URI=neo4j://localhost:7687
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=your_password
```

### 4. Run the server

```bash
go run cmd/server/main.go
```

### 5. Initialize Neo4j data (if using Neo4j)

```bash
go run cmd/neo4j/main.go
```

### 6. Test the API

```
GET http://localhost:8080/ministries
GET http://localhost:8080/departments
POST http://localhost:8080/ministries
POST http://localhost:8080/departments
GET http://localhost:8080/ministries/{id}
GET http://localhost:8080/departments/{id}
```

## 🧾 License

This project is licensed under the MIT License.
