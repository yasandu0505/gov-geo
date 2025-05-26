package repository

import "go-mysql-backend/internal/models"

// Neo4jRepo defines the interface for Neo4j repository methods.
type Neo4jRepo interface {
	GetMinistriesWithDepartments() ([]models.MinistryWithDepartments, error)
	GetMinistryByIDWithDepartments(id int) (models.MinistryWithDepartments, error)
	SeedDummyData() error
}
