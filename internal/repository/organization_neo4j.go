package repository

import (
	"fmt"
	"go-mysql-backend/internal/db"
	"go-mysql-backend/internal/models"
)

// OrganizationNeo4jRepository handles Neo4j operations for organization data
type OrganizationNeo4jRepository struct {
	Neo4j *db.Neo4jInterface
}

// NewOrganizationNeo4jRepository creates a new Neo4j repository for organizations
func NewOrganizationNeo4jRepository(neo4j *db.Neo4jInterface) *OrganizationNeo4jRepository {
	return &OrganizationNeo4jRepository{Neo4j: neo4j}
}

// GetAllDepartments retrieves all departments from Neo4j
func (r *OrganizationNeo4jRepository) GetAllDepartments() ([]models.Department, error) {
	query := `
	MATCH (d:department)
	OPTIONAL MATCH (m:ministry)-[:HAS_DEPARTMENT]->(d)
	RETURN d.id as id, d.name as name, d.google_map_script as google_map_script, 
		   CASE WHEN m IS NOT NULL THEN substring(m.id, 9) ELSE "" END as ministry_id
	`

	results, err := r.Neo4j.ExecuteQuery(query, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to get departments from Neo4j: %w", err)
	}

	departments := make([]models.Department, 0, len(results))
	for _, result := range results {
		// Convert ministry_id string to int
		var ministryID int
		if mid, ok := result["ministry_id"].(string); ok && mid != "" {
			fmt.Sscanf(mid, "%d", &ministryID)
		}

		// Parse department ID from the format "dept_X"
		var deptID int
		if did, ok := result["id"].(string); ok {
			fmt.Sscanf(did, "dept_%d", &deptID)
		}

		// Safely handle name and google_map_script
		name := ""
		if nameVal, ok := result["name"]; ok && nameVal != nil {
			if nameStr, ok := nameVal.(string); ok {
				name = nameStr
			}
		}

		mapScript := ""
		if scriptVal, ok := result["google_map_script"]; ok && scriptVal != nil {
			if scriptStr, ok := scriptVal.(string); ok {
				mapScript = scriptStr
			}
		}

		dept := models.Department{
			ID:                deptID,
			Name:              name,
			MinistryID:        ministryID,
			Google_map_script: mapScript,
		}
		departments = append(departments, dept)
	}

	return departments, nil
}

// GetMinistriesWithDepartments retrieves all ministries with their departments from Neo4j
func (r *OrganizationNeo4jRepository) GetMinistriesWithDepartments() ([]MinistryWithDepartments, error) {
	// First get all ministries
	ministryQuery := `
	MATCH (m:ministry)
	RETURN m.id as id, m.name as name, m.google_map_script as google_map_script
	`
	ministryResults, err := r.Neo4j.ExecuteQuery(ministryQuery, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to get ministries from Neo4j: %w", err)
	}

	// Create a map to store ministries by ID
	ministriesMap := make(map[string]*MinistryWithDepartments)
	for _, result := range ministryResults {
		// Parse ministry ID from the format "ministry_X"
		var ministryID int
		idStr := result["id"].(string)
		fmt.Sscanf(idStr, "ministry_%d", &ministryID)

		ministriesMap[idStr] = &MinistryWithDepartments{
			Ministry: models.Ministry{
				ID:                ministryID,
				Name:              result["name"].(string),
				Google_map_script: result["google_map_script"].(string),
			},
			Departments: []models.Department{},
		}
	}

	// Now get all departments with their ministry relationships
	deptQuery := `
	MATCH (m:ministry)-[:HAS_DEPARTMENT]->(d:department)
	RETURN m.id as ministry_id, d.id as id, d.name as name, d.google_map_script as google_map_script
	`
	deptResults, err := r.Neo4j.ExecuteQuery(deptQuery, map[string]interface{}{})
	if err != nil {
		return nil, fmt.Errorf("failed to get departments from Neo4j: %w", err)
	}

	// Add departments to their ministries
	for _, result := range deptResults {
		ministryID := result["ministry_id"].(string)

		// Parse department ID from the format "dept_X"
		var deptID int
		if did, ok := result["id"].(string); ok {
			fmt.Sscanf(did, "dept_%d", &deptID)
		}

		if ministry, ok := ministriesMap[ministryID]; ok {
			dept := models.Department{
				ID:                deptID,
				Name:              result["name"].(string),
				MinistryID:        ministry.ID,
				Google_map_script: result["google_map_script"].(string),
			}
			ministry.Departments = append(ministry.Departments, dept)
		}
	}

	// Convert map to slice
	ministries := make([]MinistryWithDepartments, 0, len(ministriesMap))
	for _, ministry := range ministriesMap {
		ministries = append(ministries, *ministry)
	}

	return ministries, nil
}

// CreateMinistry creates a ministry in Neo4j
func (r *OrganizationNeo4jRepository) CreateMinistry(ministry models.Ministry) (int, error) {
	query := `
	MATCH (g:government {id: $govID})
	MERGE (m:ministry {id: $id, name: $name})
	SET m.google_map_script = $mapScript
	MERGE (g)-[:HAS_MINISTRY]->(m)
	RETURN m
	`

	// Generate a new ID if one isn't provided
	ministryID := ministry.ID
	if ministryID == 0 {
		// Get the highest ID currently in use
		maxIDQuery := `
		MATCH (m:ministry)
		RETURN max(toInteger(substring(m.id, 9))) as max_id
		`
		results, err := r.Neo4j.ExecuteQuery(maxIDQuery, map[string]interface{}{})
		if err != nil {
			return 0, fmt.Errorf("failed to generate ministry ID: %w", err)
		}

		if len(results) > 0 && results[0]["max_id"] != nil {
			ministryID = results[0]["max_id"].(int) + 1
		} else {
			ministryID = 1
		}
	}

	params := map[string]interface{}{
		"govID":     "gov_01",
		"id":        fmt.Sprintf("ministry_%d", ministryID),
		"name":      ministry.Name,
		"mapScript": ministry.Google_map_script,
	}

	_, err := r.Neo4j.ExecuteQuery(query, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create ministry in Neo4j: %w", err)
	}

	return ministryID, nil
}

// CreateDepartment creates a department in Neo4j
func (r *OrganizationNeo4jRepository) CreateDepartment(dept models.Department) (int, error) {
	query := `
	MATCH (m:ministry {id: $ministryID})
	MERGE (d:department {id: $id, name: $name})
	SET d.google_map_script = $mapScript
	MERGE (m)-[:HAS_DEPARTMENT]->(d)
	RETURN d
	`

	// Generate a new ID if one isn't provided
	deptID := dept.ID
	if deptID == 0 {
		// Get the highest ID currently in use
		maxIDQuery := `
		MATCH (d:department)
		RETURN max(toInteger(substring(d.id, 5))) as max_id
		`
		results, err := r.Neo4j.ExecuteQuery(maxIDQuery, map[string]interface{}{})
		if err != nil {
			return 0, fmt.Errorf("failed to generate department ID: %w", err)
		}

		if len(results) > 0 && results[0]["max_id"] != nil {
			deptID = results[0]["max_id"].(int) + 1
		} else {
			deptID = 1
		}
	}

	params := map[string]interface{}{
		"ministryID": fmt.Sprintf("ministry_%d", dept.MinistryID),
		"id":         fmt.Sprintf("dept_%d", deptID),
		"name":       dept.Name,
		"mapScript":  dept.Google_map_script,
	}

	_, err := r.Neo4j.ExecuteQuery(query, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create department in Neo4j: %w", err)
	}

	return deptID, nil
}

// GetMinistryByID retrieves a ministry by ID from Neo4j
func (r *OrganizationNeo4jRepository) GetMinistryByID(id int) (models.Ministry, error) {
	query := `
	MATCH (m:ministry {id: $id})
	RETURN m.id as id, m.name as name, m.google_map_script as google_map_script
	`
	params := map[string]interface{}{
		"id": fmt.Sprintf("ministry_%d", id),
	}

	results, err := r.Neo4j.ExecuteQuery(query, params)
	if err != nil {
		return models.Ministry{}, fmt.Errorf("failed to get ministry from Neo4j: %w", err)
	}

	if len(results) == 0 {
		return models.Ministry{}, fmt.Errorf("ministry not found with id: %d", id)
	}

	// Parse ministry ID from the format "ministry_X"
	var ministryID int
	if idStr, ok := results[0]["id"].(string); ok {
		fmt.Sscanf(idStr, "ministry_%d", &ministryID)
	}

	ministry := models.Ministry{
		ID:                ministryID,
		Name:              results[0]["name"].(string),
		Google_map_script: results[0]["google_map_script"].(string),
	}

	return ministry, nil
}

// GetDepartmentByID retrieves a department by ID from Neo4j
func (r *OrganizationNeo4jRepository) GetDepartmentByID(id int) (*models.Department, error) {
	query := `
	MATCH (d:department {id: $id})
	OPTIONAL MATCH (m:ministry)-[:HAS_DEPARTMENT]->(d)
	RETURN d.id as id, d.name as name, d.google_map_script as google_map_script, 
		   CASE WHEN m IS NOT NULL THEN substring(m.id, 9) ELSE "" END as ministry_id
	`
	params := map[string]interface{}{
		"id": fmt.Sprintf("dept_%d", id),
	}

	results, err := r.Neo4j.ExecuteQuery(query, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get department from Neo4j: %w", err)
	}

	if len(results) == 0 {
		return nil, nil
	}

	// Convert ministry_id string to int
	var ministryID int
	if mid, ok := results[0]["ministry_id"].(string); ok && mid != "" {
		fmt.Sscanf(mid, "%d", &ministryID)
	}

	// Parse department ID from the format "dept_X"
	var deptID int
	if did, ok := results[0]["id"].(string); ok {
		fmt.Sscanf(did, "dept_%d", &deptID)
	}

	dept := &models.Department{
		ID:                deptID,
		Name:              results[0]["name"].(string),
		MinistryID:        ministryID,
		Google_map_script: results[0]["google_map_script"].(string),
	}

	return dept, nil
}
