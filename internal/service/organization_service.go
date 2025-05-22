package service

import (
	"go-mysql-backend/internal/models"
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

func (s *OrganizationService) GetMinistriesWithDepartmentsPaginated(limit, offset int) ([]repository.MinistryWithDepartments, error) {
	return repository.GetMinistriesWithDepartmentsPaginated(s.Repo, limit, offset)
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
	ministry, err := s.Repo.GetMinistryByID(id)
	if err != nil {
		return models.Ministry{}, err
	}
	return ministry, nil
}

func (s *OrganizationService) GetDepartmentByID(id int) (*models.Department, error) {
	return s.Repo.GetDepartmentByID(id)
}
