package admin

import (
	"testing"
	"time"

	"github.com/k/iRegistro/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockAdminRepo struct {
	mock.Mock
}

func (m *MockAdminRepo) CreateAuditLog(log *domain.AuditLog) error {
	return m.Called(log).Error(0)
}
func (m *MockAdminRepo) GetAuditLogs(schoolID *uint, limit, offset int) ([]domain.AuditLog, error) {
	args := m.Called(schoolID, limit, offset)
	return args.Get(0).([]domain.AuditLog), args.Error(1)
}
func (m *MockAdminRepo) GetSchoolSettings(schoolID uint) ([]domain.SchoolSettings, error) {
	args := m.Called(schoolID)
	return args.Get(0).([]domain.SchoolSettings), args.Error(1)
}
func (m *MockAdminRepo) UpsertSchoolSetting(setting *domain.SchoolSettings) error {
	return m.Called(setting).Error(0)
}
func (m *MockAdminRepo) CreateUserImport(imp *domain.UserImport) error {
	args := m.Called(imp)
	imp.ID = 1
	return args.Error(0)
}
func (m *MockAdminRepo) GetUserImport(id uint) (*domain.UserImport, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.UserImport), args.Error(1)
}
func (m *MockAdminRepo) UpdateUserImport(imp *domain.UserImport) error {
	return m.Called(imp).Error(0)
}
func (m *MockAdminRepo) CreateDataExport(exp *domain.DataExport) error {
	args := m.Called(exp)
	exp.ID = 1
	return args.Error(0)
}
func (m *MockAdminRepo) GetDataExport(id uint) (*domain.DataExport, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.DataExport), args.Error(1)
}
func (m *MockAdminRepo) UpdateDataExport(exp *domain.DataExport) error {
	return m.Called(exp).Error(0)
}

type MockUserRepo struct {
	mock.Mock
}

// Implement only needed methods from UserRepository interface
func (m *MockUserRepo) CreateUser(u *domain.User) error                                { return m.Called(u).Error(0) }
func (m *MockUserRepo) FindByEmail(email string) (*domain.User, error)                 { return nil, nil }
func (m *MockUserRepo) UpdateUser(u *domain.User) error                                { return nil }
func (m *MockUserRepo) FindByID(id uint) (*domain.User, error)                         { return nil, nil }
func (m *MockUserRepo) SaveRefreshToken(token *domain.RefreshToken) error              { return nil }
func (m *MockUserRepo) GetSession(token string) (*domain.Session, error)               { return nil, nil }
func (m *MockUserRepo) BlockUser(id uint) error                                        { return nil }
func (m *MockUserRepo) IncrementFailedLogin(id uint) error                             { return nil }
func (m *MockUserRepo) ResetFailedLogin(id uint) error                                 { return nil }
func (m *MockUserRepo) UpdatePassword(id uint, hash string) error                      { return nil }
func (m *MockUserRepo) SaveResetToken(token string, exp time.Time, email string) error { return nil }
func (m *MockUserRepo) ValidateResetToken(token string) (*domain.User, error)          { return nil, nil }

// --- Tests ---

func TestAuditLog(t *testing.T) {
	mockRepo := new(MockAdminRepo)
	svc := NewAuditService(mockRepo)

	mockRepo.On("CreateAuditLog", mock.MatchedBy(func(l *domain.AuditLog) bool {
		return l.Action == "LOGIN" && l.UserID == 1
	})).Return(nil)

	err := svc.LogAction(nil, 1, "LOGIN", "AUTH", "", "127.0.0.1", nil)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateSettings(t *testing.T) {
	mockRepo := new(MockAdminRepo)
	auditSvc := NewAuditService(mockRepo)
	// We need mocked repositories
	// Using nil for userRepo/schoolRepo as UpdateSetting doesn't use them in current simple impl
	svc := NewAdminService(mockRepo, nil, nil, auditSvc)

	// Expect Upsert
	mockRepo.On("UpsertSchoolSetting", mock.MatchedBy(func(s *domain.SchoolSettings) bool {
		return s.Key == "theme" && s.SchoolID == 1
	})).Return(nil)

	// Expect Audit Log
	mockRepo.On("CreateAuditLog", mock.Anything).Return(nil)

	err := svc.UpdateSchoolSetting(1, 10, "theme", map[string]interface{}{"color": "blue"})
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
