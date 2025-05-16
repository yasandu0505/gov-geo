package main

import (
	"go-mysql-backend/internal/db"
	"go-mysql-backend/routes"
	"log"
	"net/http"
)

func main() {
	db.InitMySQL()
	router := routes.SetupRouter()

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
