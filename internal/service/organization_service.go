package service

import (
	"go-mysql-backend/internal/repository"
)

type OrganizationService struct {
	Repo *repository.OrganizationRepository
}

func NewOrganizationService(repo *repository.OrganizationRepository) *OrganizationService {
	return &OrganizationService{Repo: repo}
}

func (s *OrganizationService) GetMinistriesWithDepartments() ([]repository.MinistryWithDepartments, error) {
	return repository.GetMinistriesWithDepartments(s.Repo)
}
