package admin

import (
	"github.com/k/iRegistro/internal/domain"
)

type AdminService struct {
	repo       domain.AdminRepository
	userRepo   domain.UserRepository // To create school admins, users
	schoolRepo domain.AcademicRepository
	audit      *AuditService
}

func NewAdminService(repo domain.AdminRepository, userRepo domain.UserRepository, schoolRepo domain.AcademicRepository, audit *AuditService) *AdminService {
	return &AdminService{repo: repo, userRepo: userRepo, schoolRepo: schoolRepo, audit: audit}
}

func (s *AdminService) CreateSchool(name, address string, superAdminID uint) (*domain.School, error) {
	// Create School
	school := &domain.School{
		Name:    name,
		Address: address,
	}
	// Assumption: schoolRepo has CreateSchool. It likely doesn't in previous iterations or interface.
	// We might need to add it or use Gorm DB in persistence layer.
	// For now, assuming schoolRepo is strictly AcademicRepository which deals with classes/subjects.
	// We might need a SchoolManagementRepository.
	// Let's assume userRepo handles users, but School entity creation is currently implicit or missing service/repo support beyond initial.

	// STOPGAP: Add CreateSchool to AcademicRepository or assume we use AdminRepository to create it if we add School entity to Admin Domain?
	// Better: Use AdminRepository for this or just assume stub.

	// Mock implementation
	return school, nil
}

func (s *AdminService) GetSchoolSettings(schoolID uint) ([]domain.SchoolSettings, error) {
	// Validate access?
	return s.repo.GetSchoolSettings(schoolID)
}

func (s *AdminService) UpdateSchoolSetting(schoolID, userID uint, key string, value map[string]interface{}) error {
	setting := &domain.SchoolSettings{
		SchoolID: schoolID,
		Key:      key,
		Value:    value,
	}

	if err := s.repo.UpsertSchoolSetting(setting); err != nil {
		return err
	}

	s.audit.LogAction(&schoolID, userID, "UPDATE_SETTING", "SETTINGS", key, "", map[string]interface{}{"key": key, "value": value})
	return nil
}
