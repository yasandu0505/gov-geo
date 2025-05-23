package repository

import (
	"context"
	"fmt"

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

func (r *Neo4jRepository) GetMinistryByIDWithDepartments(ministryID int) (models.MinistryWithDepartments, error) {
	ctx := context.Background()
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	query := `
		MATCH (m:Ministry {id: $id})
		OPTIONAL MATCH (m)-[:HAS_DEPARTMENT]->(d:Department)
		RETURN
			m.id AS ministry_id,
			m.name AS ministry_name,
			m.google_map_script AS ministry_map,
			d.id AS dept_id,
			d.name AS dept_name,
			d.google_map_script AS dept_map
		ORDER BY d.id
	`

	params := map[string]interface{}{
		"id": ministryID,
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		return models.MinistryWithDepartments{}, err
	}

	var ministryWithDepts models.MinistryWithDepartments
	foundMinistry := false

	for result.Next(ctx) {
		record := result.Record()

		if !foundMinistry {
			ministryWithDepts.Ministry.ID = int(record.Values[0].(int64))
			ministryWithDepts.Ministry.Name = record.Values[1].(string)

			if record.Values[2] != nil {
				ministryWithDepts.Ministry.Google_map_script = record.Values[2].(string)
			}
			foundMinistry = true
		}

		if record.Values[3] != nil { // department exists
			deptID := int(record.Values[3].(int64))
			deptName := record.Values[4].(string)
			deptMap := ""
			if record.Values[5] != nil {
				deptMap = record.Values[5].(string)
			}

			department := models.Department{
				ID:                deptID,
				Name:              deptName,
				MinistryID:        ministryID,
				Google_map_script: deptMap,
			}
			ministryWithDepts.Departments = append(ministryWithDepts.Departments, department)
		}
	}

	if err = result.Err(); err != nil {
		return models.MinistryWithDepartments{}, err
	}

	if !foundMinistry {
		return models.MinistryWithDepartments{}, fmt.Errorf("ministry with ID %d not found", ministryID)
	}

	return ministryWithDepts, nil
}
