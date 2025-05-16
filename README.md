# Go MySQL Backend API

A simple, modular RESTful API built in Go (Golang) that connects to a MySQL database. It follows a clean architecture with proper separation of concerns (handlers, services, repositories, and models).

## 📁 Project Structure

go-mysql-backend/
│
├── cmd/
│   └── server/             # Main entry point of the application
│
├── config/                 # Configuration settings (e.g. database credentials)
│
├── internal/
│   ├── db/                 # Database connection logic
│   ├── models/             # Data models (structs)
│   ├── repository/         # Data access layer (SQL queries)
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
- MySQL
- Gorilla Mux (HTTP router)
- Git

## 🚀 Getting Started

### 1. Clone the repository

git clone https://github.com/yasandu0505/gov-geo.git
cd gov-geo

### 2. Set up MySQL

CREATE DATABASE testdb;
USE testdb;

CREATE TABLE users (
  id INT AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100),
  email VARCHAR(100)
);

INSERT INTO users (name, email) VALUES
('Alice', 'alice@example.com'),
('Bob', 'bob@example.com');

### 3. Configure database credentials

Edit internal/db/mysql.go:

dsn := "root:yourpassword@tcp(127.0.0.1:3306)/testdb"

Replace `yourpassword` with your actual MySQL password.

### 4. Run the server

go run cmd/server/main.go

### 5. Test the API

GET http://localhost:8080/users

## 🧾 License

This project is licensed under the MIT License.
