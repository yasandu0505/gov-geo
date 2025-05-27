package repository

import "go-mysql-backend/internal/models"

type PostgresRepo interface {
	GetMinistriesWithDepartments() ([]models.MinistryWithDepartments, error)
	GetMinistriesWithDepartmentsPaginated(limit, offset int) ([]models.MinistryWithDepartments, error)
	GetAllDepartments() ([]models.Department, error)
	CreateMinistry(ministry models.Ministry) (int, error)
	CreateDepartment(dept models.Department) (int, error)
	GetMinistryByID(id int) (models.Ministry, error)
	GetMinistryByIDWithDepartments(id int) (models.MinistryWithDepartments, error)
	GetDepartmentByID(id int) (*models.Department, error)
}
