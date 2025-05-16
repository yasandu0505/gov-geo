package service

import (
	"go-mysql-backend/internal/models"
	"go-mysql-backend/internal/repository"
)

type UserService interface {
	GetAllUsers() ([]models.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{r}
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAll()
}
