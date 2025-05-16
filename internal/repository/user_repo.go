package repository

import (
	"database/sql"
	"log"

	"go-mysql-backend/internal/models"
)

type UserRepository interface {
	GetAll() ([]models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetAll() ([]models.User, error) {
	rows, err := r.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			log.Println("Scan error:", err)
			continue
		}
		users = append(users, u)
	}
	return users, nil
}
