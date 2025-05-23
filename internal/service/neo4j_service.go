package service

import (
	"go-mysql-backend/internal/models"
	"go-mysql-backend/internal/repository"
)

type Neo4JService struct {
	Repo *repository.Neo4jRepository
}

func NewNeo4JService(repo *repository.Neo4jRepository) *Neo4JService {
	return &Neo4JService{Repo: repo}
}

func (s *Neo4JService) GetMinistriesWithDepartments() ([]models.MinistryWithDepartments, error) {
	return s.Repo.GetMinistriesWithDepartments()
}
