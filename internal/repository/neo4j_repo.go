package repository

import (
	"context"
	"fmt"
	"math/rand"
	"time"

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
func generateRandomMinistryName() string {
	prefixes := []string{"Ministry of", "Department of", "Office of"}
	topics := []string{"Innovation", "Agriculture", "Wellbeing", "Technology", "Security", "Environment", "Commerce"}

	return fmt.Sprintf("%s %s", randomItem(prefixes), randomItem(topics))
}

func generateRandomDepartmentName() string {
	roles := []string{"Planning", "Operations", "Research", "Development", "Logistics", "Services"}
	return fmt.Sprintf("%s Department", randomItem(roles))
}

func randomItem(list []string) string {
	return list[rand.Intn(len(list))]
}

func (r *Neo4jRepository) SeedDummyData() error {
	ctx := context.Background()
	session := r.Driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	rand.Seed(time.Now().UnixNano())

	numMinistries := 200
	departmentsPerMinistry := 10

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (interface{}, error) {
		deptGlobalID := 1 // ensure unique department IDs globally

		for i := 1; i <= numMinistries; i++ {
			ministryID := i
			ministryName := generateRandomMinistryName()
			ministryMap := fmt.Sprintf("<script src='map/%d.js'></script>", ministryID)

			// Create Ministry
			_, err := tx.Run(ctx, `
				CREATE (m:Ministry {id: $id, name: $name, google_map_script: $map})
			`, map[string]interface{}{
				"id":   ministryID,
				"name": ministryName,
				"map":  ministryMap,
			})
			if err != nil {
				return nil, err
			}

			// Create Departments
			for j := 0; j < departmentsPerMinistry; j++ {
				deptID := deptGlobalID
				deptName := generateRandomDepartmentName()
				deptMap := fmt.Sprintf("<script src='dept/%d.js'></script>", deptID)

				_, err := tx.Run(ctx, `
					MATCH (m:Ministry {id: $ministryID})
					CREATE (d:Department {
						id: $deptID,
						name: $name,
						google_map_script: $map
					})
					CREATE (m)-[:HAS_DEPARTMENT]->(d)
				`, map[string]interface{}{
					"ministryID": ministryID,
					"deptID":     deptID,
					"name":       deptName,
					"map":        deptMap,
				})
				if err != nil {
					return nil, err
				}

				deptGlobalID++
			}
		}
		return nil, nil
	})

	return err
}
