package service_test

import (
	"database/sql"
	"testing"

	"go-mysql-backend/internal/models"
	"go-mysql-backend/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPostgresRepo is a mock implementation of PostgresRepo interface
type MockPostgresRepo struct {
	mock.Mock
}

func (m *MockPostgresRepo) GetMinistriesWithDepartments() ([]models.MinistryWithDepartments, error) {
	args := m.Called()
	return args.Get(0).([]models.MinistryWithDepartments), args.Error(1)
}

func (m *MockPostgresRepo) GetMinistriesWithDepartmentsPaginated(limit, offset int) ([]models.MinistryWithDepartments, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]models.MinistryWithDepartments), args.Error(1)
}

func (m *MockPostgresRepo) GetAllDepartments() ([]models.Department, error) {
	args := m.Called()
	return args.Get(0).([]models.Department), args.Error(1)
}

func (m *MockPostgresRepo) CreateMinistry(ministry models.Ministry) (int, error) {
	args := m.Called(ministry)
	return args.Int(0), args.Error(1)
}

func (m *MockPostgresRepo) CreateDepartment(dept models.Department) (int, error) {
	args := m.Called(dept)
	return args.Int(0), args.Error(1)
}

func (m *MockPostgresRepo) GetMinistryByID(id int) (models.Ministry, error) {
	args := m.Called(id)
	return args.Get(0).(models.Ministry), args.Error(1)
}

func (m *MockPostgresRepo) GetMinistryByIDWithDepartments(id int) (models.MinistryWithDepartments, error) {
	args := m.Called(id)
	return args.Get(0).(models.MinistryWithDepartments), args.Error(1)
}

func (m *MockPostgresRepo) GetDepartmentByID(id int) (*models.Department, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Department), args.Error(1)
}

func TestPostgresGetMinistriesWithDepartments(t *testing.T) {
	mockRepo := new(MockPostgresRepo)
	service := service.NewOrganizationService(mockRepo)

	expectedMinistries := []models.MinistryWithDepartments{
		{
			Ministry: models.Ministry{
				ID:                1,
				Name:              "Ministry of Education",
				Google_map_script: "<script>map1</script>",
			},
			Departments: []models.Department{
				{
					ID:                1,
					Name:              "Primary Education",
					MinistryID:        1,
					Google_map_script: "<script>map2</script>",
				},
			},
		},
	}

	mockRepo.On("GetMinistriesWithDepartments").Return(expectedMinistries, nil)

	result, err := service.GetMinistriesWithDepartments()

	assert.NoError(t, err)
	assert.Equal(t, expectedMinistries, result)
	mockRepo.AssertExpectations(t)
}

func TestPostgresGetMinistriesWithDepartmentsPaginated(t *testing.T) {
	mockRepo := new(MockPostgresRepo)
	service := service.NewOrganizationService(mockRepo)

	limit, offset := 10, 0
	expectedMinistries := []models.MinistryWithDepartments{
		{
			Ministry: models.Ministry{
				ID:   1,
				Name: "Ministry of Education",
			},
			Departments: []models.Department{
				{
					ID:         1,
					Name:       "Primary Education",
					MinistryID: 1,
				},
			},
		},
	}

	mockRepo.On("GetMinistriesWithDepartmentsPaginated", limit, offset).Return(expectedMinistries, nil)

	result, err := service.GetMinistriesWithDepartmentsPaginated(limit, offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedMinistries, result)
	mockRepo.AssertExpectations(t)
}

func TestPostgresCreateMinistry(t *testing.T) {
	mockRepo := new(MockPostgresRepo)
	service := service.NewOrganizationService(mockRepo)

	ministry := models.Ministry{
		Name:              "New Ministry",
		Google_map_script: "<script>map</script>",
	}
	expectedID := 1

	mockRepo.On("CreateMinistry", ministry).Return(expectedID, nil)

	id, err := service.CreateMinistry(ministry)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
	mockRepo.AssertExpectations(t)
}

func TestPostgresCreateDepartment(t *testing.T) {
	mockRepo := new(MockPostgresRepo)
	service := service.NewOrganizationService(mockRepo)

	department := models.Department{
		Name:              "New Department",
		MinistryID:        1,
		Google_map_script: "<script>map</script>",
	}
	expectedID := 1

	mockRepo.On("CreateDepartment", department).Return(expectedID, nil)

	id, err := service.CreateDepartment(department)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, id)
	mockRepo.AssertExpectations(t)
}

func TestPostgresGetMinistryByID(t *testing.T) {
	mockRepo := new(MockPostgresRepo)
	service := service.NewOrganizationService(mockRepo)

	expectedMinistry := models.Ministry{
		ID:                1,
		Name:              "Ministry of Education",
		Google_map_script: "<script>map</script>",
	}

	mockRepo.On("GetMinistryByID", 1).Return(expectedMinistry, nil)

	result, err := service.GetMinistryByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedMinistry, result)
	mockRepo.AssertExpectations(t)
}

func TestPostgresGetMinistryByID_NotFound(t *testing.T) {
	mockRepo := new(MockPostgresRepo)
	service := service.NewOrganizationService(mockRepo)

	mockRepo.On("GetMinistryByID", 999).Return(models.Ministry{}, sql.ErrNoRows)

	result, err := service.GetMinistryByID(999)

	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
	assert.Equal(t, models.Ministry{}, result)
	mockRepo.AssertExpectations(t)
}

func TestPostgresGetMinistryByIDWithDepartments(t *testing.T) {
	mockRepo := new(MockPostgresRepo)
	service := service.NewOrganizationService(mockRepo)

	expectedMinistry := models.MinistryWithDepartments{
		Ministry: models.Ministry{
			ID:                1,
			Name:              "Ministry of Education",
			Google_map_script: "<script>map1</script>",
		},
		Departments: []models.Department{
			{
				ID:                1,
				Name:              "Primary Education",
				MinistryID:        1,
				Google_map_script: "<script>map2</script>",
			},
		},
	}

	mockRepo.On("GetMinistryByIDWithDepartments", 1).Return(expectedMinistry, nil)

	result, err := service.GetMinistryByIDWithDepartments(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedMinistry, result)
	mockRepo.AssertExpectations(t)
}

func TestPostgresGetDepartmentByID(t *testing.T) {
	mockRepo := new(MockPostgresRepo)
	service := service.NewOrganizationService(mockRepo)

	expectedDepartment := &models.Department{
		ID:                1,
		Name:              "Primary Education",
		MinistryID:        1,
		Google_map_script: "<script>map</script>",
	}

	mockRepo.On("GetDepartmentByID", 1).Return(expectedDepartment, nil)

	result, err := service.GetDepartmentByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedDepartment, result)
	mockRepo.AssertExpectations(t)
}

func TestPostgresGetDepartmentByID_NotFound(t *testing.T) {
	mockRepo := new(MockPostgresRepo)
	service := service.NewOrganizationService(mockRepo)

	mockRepo.On("GetDepartmentByID", 999).Return(nil, nil)

	result, err := service.GetDepartmentByID(999)

	assert.NoError(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

func TestPostgresGetAllDepartments(t *testing.T) {
	mockRepo := new(MockPostgresRepo)
	service := service.NewOrganizationService(mockRepo)

	expectedDepartments := []models.Department{
		{
			ID:                1,
			Name:              "Primary Education",
			MinistryID:        1,
			Google_map_script: "<script>map1</script>",
		},
		{
			ID:                2,
			Name:              "Secondary Education",
			MinistryID:        1,
			Google_map_script: "<script>map2</script>",
		},
	}

	mockRepo.On("GetAllDepartments").Return(expectedDepartments, nil)

	result, err := service.GetAllDepartments()

	assert.NoError(t, err)
	assert.Equal(t, expectedDepartments, result)
	mockRepo.AssertExpectations(t)
}
