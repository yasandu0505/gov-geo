package service_test

import (
	"testing"

	"go-mysql-backend/internal/models"
	"go-mysql-backend/internal/repository/mocks"
	"go-mysql-backend/internal/service"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetMinistriesWithDepartments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNeo4jRepo(ctrl)

	// Define expected output
	expected_list := []models.MinistryWithDepartments{
		{
			Ministry: models.Ministry{
				ID:   1,
				Name: "Ministry of Testing",
			},
			Departments: []models.Department{
				{
					ID:         101,
					Name:       "QA Dept",
					MinistryID: 1,
				},
			},
		},
	}

	// Setup expectations
	mockRepo.EXPECT().GetMinistriesWithDepartments().Return(expected_list, nil)

	// Inject mock into service
	s := service.NewNeo4JService(mockRepo)

	// Run the test
	result, err := s.GetMinistriesWithDepartments()

	assert.NoError(t, err)
	assert.Equal(t, expected_list, result)
}

func TestGetMinistryByIdWithDepartments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNeo4jRepo(ctrl)

	// Define expected output
	expected_ministry := models.MinistryWithDepartments{

		Ministry: models.Ministry{
			ID:   1,
			Name: "Ministry of Testing",
		},
		Departments: []models.Department{
			{
				ID:         101,
				Name:       "QA Dept",
				MinistryID: 1,
			},
		},
	}

	// Setup expectations
	mockRepo.EXPECT().GetMinistryByIDWithDepartments(1).Return(expected_ministry, nil)

	// Inject mock into service
	s := service.NewNeo4JService(mockRepo)

	// Run the test
	result, err := s.GetMinistryByIDWithDepartments(1)

	assert.NoError(t, err)
	assert.Equal(t, expected_ministry, result)
}
