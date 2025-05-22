package service

import (
	"go-mysql-backend/internal/models"
	"go-mysql-backend/internal/repository"
)

type OrganizationService struct {
	Repo repository.OrganizationRepositoryInterface
}

func NewOrganizationService(repo repository.OrganizationRepositoryInterface) *OrganizationService {
	return &OrganizationService{Repo: repo}
}

func (s *OrganizationService) GetMinistriesWithDepartments() ([]repository.MinistryWithDepartments, error) {
	return s.Repo.GetMinistriesWithDepartments()
}

func (s *OrganizationService) CreateMinistry(ministry models.Ministry) (int, error) {
	return s.Repo.CreateMinistry(ministry)
}

func (s *OrganizationService) CreateDepartment(department models.Department) (int, error) {
	return s.Repo.CreateDepartment(department)
}
func (s *OrganizationService) GetAllDepartments() ([]models.Department, error) {
	return s.Repo.GetAllDepartments()
}
func (s *OrganizationService) GetMinistryByID(id int) (models.Ministry, error) {
	return s.Repo.GetMinistryByID(id)
}

func (s *OrganizationService) GetDepartmentByID(id int) (*models.Department, error) {
	return s.Repo.GetDepartmentByID(id)
}
