package repository

import (
	"database/sql"
	"go-mysql-backend/internal/models"

	_ "github.com/lib/pq"
)

type OrganizationRepository struct {
	DB *sql.DB
}

func NewOrganizationRepository(db *sql.DB) *OrganizationRepository {
	return &OrganizationRepository{DB: db}
}

type MinistryWithDepartments struct {
	models.Ministry
	Departments []models.Department
}

func GetMinistriesWithDepartments(r *OrganizationRepository) ([]MinistryWithDepartments, error) {
	rows, err := r.DB.Query(`
        SELECT m.id, m.name, d.id, d.name, d.ministry_id
        FROM ministry m
        LEFT JOIN department d ON m.id = d.ministry_id
        ORDER BY m.id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ministriesMap := make(map[int]*MinistryWithDepartments)

	for rows.Next() {
		var mID int
		var mName string
		var dID sql.NullInt64
		var dName sql.NullString
		var dMinistryID sql.NullInt64

		err := rows.Scan(&mID, &mName, &dID, &dName, &dMinistryID)
		if err != nil {
			return nil, err
		}

		if _, exists := ministriesMap[mID]; !exists {
			ministriesMap[mID] = &MinistryWithDepartments{
				Ministry: models.Ministry{ID: mID, Name: mName},
			}
		}

		if dID.Valid {
			dept := models.Department{
				ID:         int(dID.Int64),
				Name:       dName.String,
				MinistryID: int(dMinistryID.Int64),
			}
			ministriesMap[mID].Departments = append(ministriesMap[mID].Departments, dept)
		}
	}

	// Convert map to slice
	var ministries []MinistryWithDepartments
	for _, m := range ministriesMap {
		ministries = append(ministries, *m)
	}

	return ministries, nil
}
func (r *OrganizationRepository) CreateMinistry(ministry models.Ministry) (int, error) {
	var id int
	err := r.DB.QueryRow(`INSERT INTO ministry (name) VALUES ($1) RETURNING id`, ministry.Name).Scan(&id)
	return id, err
}

func (r *OrganizationRepository) CreateDepartment(dept models.Department) (int, error) {
	var id int
	err := r.DB.QueryRow(`INSERT INTO department (name, ministry_id) VALUES ($1, $2) RETURNING id`, dept.Name, dept.MinistryID).Scan(&id)
	return id, err
}
