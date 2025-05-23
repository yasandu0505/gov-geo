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

func GetMinistriesWithDepartments(r *OrganizationRepository) ([]models.MinistryWithDepartments, error) {
	rows, err := r.DB.Query(`
        SELECT 
            m.id, m.name, m.google_map_script,
            d.id, d.name, d.ministry_id, d.google_map_script
        FROM ministry m
        LEFT JOIN department d ON m.id = d.ministry_id
        ORDER BY m.id
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ministriesMap := make(map[int]*models.MinistryWithDepartments)

	for rows.Next() {
		var mID int
		var mName string
		var mMapScript sql.NullString
		var dID sql.NullInt64
		var dName sql.NullString
		var dMinistryID sql.NullInt64
		var dMapScript sql.NullString

		err := rows.Scan(
			&mID, &mName, &mMapScript,
			&dID, &dName, &dMinistryID, &dMapScript,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := ministriesMap[mID]; !exists {
			ministriesMap[mID] = &models.MinistryWithDepartments{
				Ministry: models.Ministry{
					ID:                mID,
					Name:              mName,
					Google_map_script: mMapScript.String,
				},
			}
		}

		if dID.Valid {
			dept := models.Department{
				ID:                int(dID.Int64),
				Name:              dName.String,
				MinistryID:        int(dMinistryID.Int64),
				Google_map_script: dMapScript.String,
			}
			ministriesMap[mID].Departments = append(ministriesMap[mID].Departments, dept)
		}
	}

	var ministries []models.MinistryWithDepartments
	for _, m := range ministriesMap {
		ministries = append(ministries, *m)
	}

	return ministries, nil
}

func (r *OrganizationRepository) GetMinistriesWithDepartmentsPaginated(limit, offset int) ([]models.MinistryWithDepartments, error) {
	query := `
        SELECT 
            m.id, m.name, m.google_map_script,
            d.id, d.name, d.google_map_script, d.ministry_id
        FROM ministry m
        LEFT JOIN department d ON m.id = d.ministry_id
        ORDER BY m.id
        LIMIT $1 OFFSET $2
    `
	rows, err := r.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ministriesMap := make(map[int]*models.MinistryWithDepartments)

	for rows.Next() {
		var mID int
		var mName, mScript string
		var dID sql.NullInt64
		var dName, dScript sql.NullString
		var dMinistryID sql.NullInt64

		err := rows.Scan(&mID, &mName, &mScript, &dID, &dName, &dScript, &dMinistryID)
		if err != nil {
			return nil, err
		}

		if _, exists := ministriesMap[mID]; !exists {
			ministriesMap[mID] = &models.MinistryWithDepartments{
				Ministry: models.Ministry{
					ID:                mID,
					Name:              mName,
					Google_map_script: mScript,
				},
			}
		}

		if dID.Valid {
			dept := models.Department{
				ID:                int(dID.Int64),
				Name:              dName.String,
				Google_map_script: dScript.String,
				MinistryID:        int(dMinistryID.Int64),
			}
			ministriesMap[mID].Departments = append(ministriesMap[mID].Departments, dept)
		}
	}

	var ministries []models.MinistryWithDepartments
	for _, m := range ministriesMap {
		ministries = append(ministries, *m)
	}

	return ministries, nil
}

func (r *OrganizationRepository) GetAllDepartments() ([]models.Department, error) {
	rows, err := r.DB.Query(`SELECT id, name, ministry_id, google_map_script FROM department`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departments []models.Department
	for rows.Next() {
		var d models.Department
		err := rows.Scan(&d.ID, &d.Name, &d.MinistryID, &d.Google_map_script)
		if err != nil {
			return nil, err
		}
		departments = append(departments, d)
	}

	return departments, nil
}

func (r *OrganizationRepository) CreateMinistry(ministry models.Ministry) (int, error) {
	var id int
	err := r.DB.QueryRow(`INSERT INTO ministry (name, google_map_script) VALUES ($1, $2) RETURNING id`,
		ministry.Name, ministry.Google_map_script).Scan(&id)
	return id, err
}

func (r *OrganizationRepository) CreateDepartment(dept models.Department) (int, error) {
	var id int
	err := r.DB.QueryRow(`INSERT INTO department (name, ministry_id, google_map_script) VALUES ($1, $2, $3) RETURNING id`, dept.Name, dept.MinistryID, dept.Google_map_script).Scan(&id)
	return id, err
}

func (r *OrganizationRepository) GetMinistryByID(id int) (models.Ministry, error) {
	var ministry models.Ministry
	err := r.DB.QueryRow(`SELECT id, name, google_map_script FROM ministry WHERE id = $1`, id).Scan(
		&ministry.ID, &ministry.Name, &ministry.Google_map_script)
	if err != nil {
		return ministry, err
	}
	return ministry, nil
}

func (r *OrganizationRepository) GetDepartmentByID(id int) (*models.Department, error) {
	row := r.DB.QueryRow(`SELECT id, name, google_map_script, ministry_id FROM department WHERE id = $1`, id)

	var dept models.Department
	err := row.Scan(&dept.ID, &dept.Name, &dept.Google_map_script, &dept.MinistryID)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &dept, nil
}
