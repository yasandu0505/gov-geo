package service

import (
	"go-mysql-backend/internal/models"
	"go-mysql-backend/internal/repository"
)

func FetchUsers() ([]models.User, error) {
	return repository.GetAllUsers()
}
