package repository

import (
	"context"

	"go-mysql-backend/internal/models"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4jRepository struct {
	Driver neo4j.DriverWithContext
}

func NewNeo4jRepository(driver neo4j.DriverWithContext) *Neo4jRepository {
	return &Neo4jRepository{Driver: driver}
}

func (r *Neo4jRepository) GetMinistriesWithDepartments() ([]models.MinistryWithDepartments, error) {
	ctx := context.Background()
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := `
		MATCH (m:Ministry)-[:HAS_DEPARTMENT]->(d:Department)
		RETURN
			m.id AS ministry_id,
			m.name AS ministry_name,
			m.google_map_script AS ministry_map,
			d.id AS dept_id,
			d.name AS dept_name,
			d.google_map_script AS dept_map
		ORDER BY m.id
	`

	result, err := session.Run(ctx, query, nil)
	if err != nil {
		return nil, err
	}

	ministryMap := make(map[int]*models.MinistryWithDepartments)

	for result.Next(ctx) {
		record := result.Record()

		ministryID := int(record.Values[0].(int64)) // Correct type
		ministryName := record.Values[1].(string)

		ministryMapScript := ""
		if record.Values[2] != nil {
			ministryMapScript = record.Values[2].(string)
		}

		deptID := int(record.Values[3].(int64))
		deptName := record.Values[4].(string)

		deptMap := ""
		if record.Values[5] != nil {
			deptMap = record.Values[5].(string)
		}

		if _, exists := ministryMap[ministryID]; !exists {
			ministryMap[ministryID] = &models.MinistryWithDepartments{
				Ministry: models.Ministry{
					ID:                ministryID,
					Name:              ministryName,
					Google_map_script: ministryMapScript,
				},
			}
		}

		department := models.Department{
			ID:                deptID,
			Name:              deptName,
			MinistryID:        ministryID,
			Google_map_script: deptMap,
		}

		ministryMap[ministryID].Departments = append(ministryMap[ministryID].Departments, department)
	}

	if err = result.Err(); err != nil {
		return nil, err
	}

	var ministries []models.MinistryWithDepartments
	for _, m := range ministryMap {
		ministries = append(ministries, *m)
	}

	return ministries, nil
}
