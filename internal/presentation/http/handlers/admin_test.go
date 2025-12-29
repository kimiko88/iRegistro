package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/k/iRegistro/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAdminService struct {
	mock.Mock
}

func (m *MockAdminService) CreateUser(schoolID uint, user *domain.User, subjectIDs []uint) error {
	args := m.Called(schoolID, user, subjectIDs)
	return args.Error(0)
}

func (m *MockAdminService) GetUsers(schoolID uint) ([]domain.User, error) {
	args := m.Called(schoolID)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockAdminService) UpdateUser(id uint, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockAdminService) CreateSchool(data map[string]interface{}) (*domain.School, error) {
	args := m.Called(data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.School), args.Error(1)
}

func (m *MockAdminService) GetSchoolSettings(schoolID uint) ([]domain.SchoolSettings, error) {
	args := m.Called(schoolID)
	return args.Get(0).([]domain.SchoolSettings), args.Error(1)
}

func (m *MockAdminService) UpdateSchoolSetting(schoolID, userID uint, key string, value map[string]interface{}) error {
	args := m.Called(schoolID, userID, key, value)
	return args.Error(0)
}

func TestAdminHandler_CreateUser_SecretaryRoleValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		requesterRole domain.Role
		targetRole    string
		expectStatus  int
		expectError   bool
	}{
		{
			name:          "Secretary creates Teacher - allowed",
			requesterRole: domain.RoleSecretary,
			targetRole:    "Teacher",
			expectStatus:  http.StatusCreated,
			expectError:   false,
		},
		{
			name:          "Secretary creates Student - allowed",
			requesterRole: domain.RoleSecretary,
			targetRole:    "Student",
			expectStatus:  http.StatusCreated,
			expectError:   false,
		},
		{
			name:          "Secretary creates Admin - forbidden",
			requesterRole: domain.RoleSecretary,
			targetRole:    "Admin",
			expectStatus:  http.StatusForbidden,
			expectError:   true,
		},
		{
			name:          "Admin creates Secretary - allowed",
			requesterRole: domain.RoleAdmin,
			targetRole:    "Secretary",
			expectStatus:  http.StatusCreated,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockAdminService)
			handler := &AdminHandler{
				adminService: mockService,
			}

			// Setup mock expectation only if we expect success
			if !tt.expectError {
				mockService.On("CreateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Setup request
			reqBody := map[string]interface{}{
				"firstName": "Test",
				"lastName":  "User",
				"email":     "test@example.com",
				"role":      tt.targetRole,
				"password":  "password123",
			}
			jsonData, _ := json.Marshal(reqBody)
			c.Request = httptest.NewRequest("POST", "/admin/users", bytes.NewBuffer(jsonData))
			c.Request.Header.Set("Content-Type", "application/json")

			// Set context values
			c.Set("role", tt.requesterRole)
			c.Set("schoolID", uint(1))
			c.Set("userID", uint(100))

			// Call handler
			handler.CreateUser(c)

			// Assertions
			assert.Equal(t, tt.expectStatus, w.Code)

			if tt.expectError {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response, "error")
			} else {
				mockService.AssertExpectations(t)
			}
		})
	}
}
