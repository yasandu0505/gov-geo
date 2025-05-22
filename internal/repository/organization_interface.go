package repository

import "go-mysql-backend/internal/models"

// OrganizationRepositoryInterface defines the methods that all organization repositories must implement
type OrganizationRepositoryInterface interface {
	GetAllDepartments() ([]models.Department, error)
	GetMinistriesWithDepartments() ([]MinistryWithDepartments, error)
	CreateMinistry(ministry models.Ministry) (int, error)
	CreateDepartment(dept models.Department) (int, error)
	GetMinistryByID(id int) (models.Ministry, error)
	GetDepartmentByID(id int) (*models.Department, error)
}
