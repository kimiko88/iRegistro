package admin

import (
	"context"
	"testing"

	"github.com/k/iRegistro/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindAll(schoolID uint) ([]domain.User, error) {
	args := m.Called(schoolID)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByExternalID(ctx context.Context, externalID string) (*domain.User, error) {
	args := m.Called(ctx, externalID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

type MockAdminRepository struct {
	mock.Mock
}

// Audit
func (m *MockAdminRepository) CreateAuditLog(log *domain.AuditLog) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockAdminRepository) GetAuditLogs(schoolID *uint, limit, offset int) ([]domain.AuditLog, error) {
	args := m.Called(schoolID, limit, offset)
	return args.Get(0).([]domain.AuditLog), args.Error(1)
}

// Settings
func (m *MockAdminRepository) GetSchoolSettings(schoolID uint) ([]domain.SchoolSettings, error) {
	args := m.Called(schoolID)
	return args.Get(0).([]domain.SchoolSettings), args.Error(1)
}

func (m *MockAdminRepository) UpsertSchoolSetting(setting *domain.SchoolSettings) error {
	args := m.Called(setting)
	return args.Error(0)
}

// Imports
func (m *MockAdminRepository) CreateUserImport(imp *domain.UserImport) error {
	args := m.Called(imp)
	return args.Error(0)
}

func (m *MockAdminRepository) GetUserImport(id uint) (*domain.UserImport, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserImport), args.Error(1)
}

func (m *MockAdminRepository) UpdateUserImport(imp *domain.UserImport) error {
	args := m.Called(imp)
	return args.Error(0)
}

// Exports
func (m *MockAdminRepository) CreateDataExport(exp *domain.DataExport) error {
	args := m.Called(exp)
	return args.Error(0)
}

func (m *MockAdminRepository) GetDataExport(id uint) (*domain.DataExport, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DataExport), args.Error(1)
}

func (m *MockAdminRepository) UpdateDataExport(exp *domain.DataExport) error {
	args := m.Called(exp)
	return args.Error(0)
}

// --- Tests ---

func TestAdminService_GetUsers(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAdminRepo := new(MockAdminRepository)
	service := NewAdminService(mockAdminRepo, mockUserRepo, nil, nil)

	schoolID := uint(1)
	expectedUsers := []domain.User{
		{ID: 1, Email: "test@example.com", Role: domain.RoleTeacher},
		{ID: 2, Email: "student@example.com", Role: domain.RoleStudent},
	}

	mockUserRepo.On("FindAll", schoolID).Return(expectedUsers, nil)

	users, err := service.GetUsers(schoolID)

	assert.NoError(t, err)
	assert.Equal(t, len(expectedUsers), len(users))
	assert.Equal(t, expectedUsers[0].Email, users[0].Email)
	mockUserRepo.AssertExpectations(t)
}

func TestAdminService_CreateUser(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAdminRepo := new(MockAdminRepository)
	service := NewAdminService(mockAdminRepo, mockUserRepo, nil, nil)

	schoolID := uint(1)
	user := &domain.User{
		Email: "new@example.com",
		Role:  domain.RoleTeacher,
	}

	mockUserRepo.On("Create", user).Return(nil)

	err := service.CreateUser(schoolID, user)

	assert.NoError(t, err)
	assert.Equal(t, schoolID, user.SchoolID)
	mockUserRepo.AssertExpectations(t)
}

// Disable UpdateSchoolSetting test for now if AuditService logic is complex or uses different repo interfaces.
// But based on signature it should work if we mocked mockAdminRepo correctly.
func TestAdminService_UpdateSchoolSetting(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockAdminRepo := new(MockAdminRepository)
	mockAuditService := NewAuditService(mockAdminRepo)

	service := NewAdminService(mockAdminRepo, mockUserRepo, nil, mockAuditService)

	schoolID := uint(1)
	userID := uint(100)
	key := "theme"
	value := map[string]interface{}{"color": "dark"}

	// Expect UpsertSchoolSetting
	mockAdminRepo.On("UpsertSchoolSetting", mock.MatchedBy(func(s *domain.SchoolSettings) bool {
		return s.SchoolID == schoolID && s.Key == key
	})).Return(nil)

	// Expect CreateAuditLog
	mockAdminRepo.On("CreateAuditLog", mock.MatchedBy(func(log *domain.AuditLog) bool {
		return log.Action == "UPDATE_SETTING" && log.UserID == userID
	})).Return(nil)

	err := service.UpdateSchoolSetting(schoolID, userID, key, value)

	assert.NoError(t, err)
	mockAdminRepo.AssertExpectations(t)
}
