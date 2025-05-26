package service

import (
	"go-mysql-backend/internal/models"
	"go-mysql-backend/internal/repository"
)

type Neo4JService struct {
	Repo repository.Neo4jRepo // Use the interface here!
}

func NewNeo4JService(repo repository.Neo4jRepo) *Neo4JService {
	return &Neo4JService{Repo: repo}
}

func (s *Neo4JService) GetMinistriesWithDepartments() ([]models.MinistryWithDepartments, error) {
	return s.Repo.GetMinistriesWithDepartments()
}

func (s *Neo4JService) GetMinistryByIDWithDepartments(id int) (models.MinistryWithDepartments, error) {
	return s.Repo.GetMinistryByIDWithDepartments(id)
}

func (s *Neo4JService) SeedDummyData() error {
	return s.Repo.SeedDummyData()
}
