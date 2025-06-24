# Government Geographic Data API 🌍

A modern RESTful API built with Go that manages government organizations (ministries and departments) with their geographic data. The API supports both PostgreSQL and Neo4j databases, offering flexibility in data storage solutions.

## 🌟 Features

- **Dual Database Support**: Choose between PostgreSQL or Neo4j
- **Ministry Management**: Create and manage government ministries
- **Department Management**: Handle departments within ministries
- **Geographic Data**: Store and retrieve geographic information
- **RESTful API**: Clean and well-documented endpoints
- **Pagination**: Efficient data retrieval with pagination support
- **CORS Support**: Secure cross-origin requests
- **Clean Architecture**: Well-organized code structure following best practices

## 🛠️ Tech Stack

- **Go (Golang)** v1.24.3
- **PostgreSQL** - Primary relational database
- **Neo4j** - Graph database option
- **Gorilla Mux** - HTTP router and URL matcher
- **Testify** - Testing toolkit
- **CORS** - Cross-Origin Resource Sharing support
- **Godotenv** - Environment configuration

## 📁 Project Structure

```
gov-geo/
├── cmd/
│   └── server/             # Application entry point
├── config/                 # Configuration settings
├── internal/
│   ├── db/                # Database connections
│   ├── models/            # Data structures
│   ├── repository/        # Data access layer
│   ├── service/           # Business logic
│   ├── handlers/          # HTTP handlers
│   └── errors/            # Custom error definitions
├── routes/                # API route definitions
└── tests/                 # Test suites
```

## 🚀 Getting Started

### Prerequisites

- Go 1.24.3 or later
- PostgreSQL 14+ or Neo4j 4.4+
- Git

### Installation

1. **Clone the Repository**
   ```bash
   git clone https://github.com/yasandu0505/gov-geo.git
   cd gov-geo
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Database Setup**

   **For PostgreSQL:**
   ```sql
   -- Connect to PostgreSQL
   psql -U postgres

   -- Create Database
   CREATE DATABASE gov_geo;

   -- Connect to the new database
   \c gov_geo

   -- Create required tables (run these in order)
   CREATE TABLE ministry (
       id SERIAL PRIMARY KEY,
       name VARCHAR(255) NOT NULL,
       google_map_script TEXT
   );

   CREATE TABLE department (
       id SERIAL PRIMARY KEY,
       name VARCHAR(255) NOT NULL,
       ministry_id INTEGER REFERENCES ministry(id),
       google_map_script TEXT
   );
   ```

   **For Neo4j:**
   - Install Neo4j Desktop from [Neo4j Download Page](https://neo4j.com/download/)
   - Create a new project and database
   - Set password for the database
   - Start the database
   - The default URI will be: `bolt://localhost:7687`

4. **Environment Configuration**

   Create a `.env` file in the root directory:
   ```env
   # Server Configuration
   PORT=8080
   ENV=development

   # Database Selection (choose one)
   DATABASE_TYPE=postgres    # Use 'postgres' for PostgreSQL
   # DATABASE_TYPE=neo4j    # Use 'neo4j' for Neo4j

   # PostgreSQL Configuration
    DATABASE_URL=     #connection String 

   # Neo4j Configuration
   NEO4J_URI=bolt://localhost:7687  # Default Neo4j URI
   NEO4J_USER=neo4j            # Default username
   NEO4J_PASSWORD=your_password    # Your Neo4j password

   # CORS Configuration (for development)
   CORS_ALLOWED_ORIGINS=http://localhost:5173
   ```

   > ⚠️ **Important Notes**: 
   > - Never commit your `.env` file to version control
   > - Keep your database passwords secure
   > - For production, use strong passwords and different port numbers
   > - Make sure to set appropriate CORS origins in production

5. **Run the Application**
   ```bash
   # Development mode
   go run cmd/server/main.go

   # Or build and run
   go build -o gov-geo cmd/server/main.go
   ./gov-geo
   ```

   The server will start on `http://localhost:8080` (or the port specified in your .env file)

## 📡 API Endpoints

### Ministries

| Method | Endpoint | Description | Request Body Example |
|--------|----------|-------------|---------------------|
| GET | `/ministries` | Get all ministries with departments | - |
| GET | `/ministries/paginated?limit=10&offset=0` | Get paginated ministries | - |
| GET | `/ministries/{id}` | Get ministry by ID | - |
| POST | `/ministries` | Create new ministry | `{"name": "Ministry of Education", "google_map_script": "<script>...</script>"}` |

### Departments

| Method | Endpoint | Description | Request Body Example |
|--------|----------|-------------|---------------------|
| GET | `/departments` | Get all departments | - |
| GET | `/departments/{id}` | Get department by ID | - |
| POST | `/departments` | Create new department | `{"name": "Primary Education", "ministry_id": 1, "google_map_script": "<script>...</script>"}` |

## 🧪 Testing

Run all tests:
```bash
go test -v ./...
```

Run specific test suite:
```bash
# Run service tests
go test -v ./internal/service/service_tests/

# Run with coverage
go test -v -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🔒 Security Features

- Environment-based configuration
- CORS protection with configurable origins
- Input validation for all endpoints
- SQL injection protection
- Error handling with custom error types
- No hardcoded credentials
- Secure password handling

## 🛠️ Troubleshooting

Common issues and solutions:

1. **Database Connection Failed**
   - Check if database service is running
   - Verify credentials in `.env` file
   - Ensure database exists and is accessible

2. **CORS Errors**
   - Check if `CORS_ALLOWED_ORIGINS` matches your frontend URL
   - Verify HTTP/HTTPS protocol matches

3. **Build Errors**
   - Run `go mod tidy` to clean up dependencies
   - Ensure Go version matches requirement (1.24.3+)

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support:
- Open an issue in the GitHub repository
- Contact the maintainers
- Check the [Troubleshooting](#troubleshooting) section

## 🙏 Acknowledgments

- Thanks to all contributors who have helped shape this project
- Special thanks to the Go community for excellent documentation and tools
- Built with ❤️ by the gov-geo team

