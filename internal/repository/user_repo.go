package repository

import (
	"go-mysql-backend/internal/db"
	"go-mysql-backend/internal/models"
)

func GetAllUsers() ([]models.User, error) {
	rows, err := db.DB.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Name, &u.Email)
		users = append(users, u)
	}
	return users, nil
}
