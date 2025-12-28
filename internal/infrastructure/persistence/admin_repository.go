package persistence

import (
	"github.com/k/iRegistro/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// --- Audit ---

func (r *AdminRepository) CreateAuditLog(log *domain.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *AdminRepository) GetAuditLogs(schoolID *uint, limit, offset int) ([]domain.AuditLog, error) {
	var logs []domain.AuditLog
	query := r.db.Order("timestamp desc").Limit(limit).Offset(offset)

	if schoolID != nil {
		query = query.Where("school_id = ?", *schoolID)
	} else {
		query = query.Where("school_id IS NULL") // SuperAdmin logs
	}

	err := query.Find(&logs).Error
	return logs, err
}

// --- Settings ---

func (r *AdminRepository) GetSchoolSettings(schoolID uint) ([]domain.SchoolSettings, error) {
	var settings []domain.SchoolSettings
	err := r.db.Where("school_id = ?", schoolID).Find(&settings).Error
	return settings, err
}

func (r *AdminRepository) UpsertSchoolSetting(setting *domain.SchoolSettings) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "school_id"}, {Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(setting).Error
}

// --- Imports ---

func (r *AdminRepository) CreateUserImport(imp *domain.UserImport) error {
	return r.db.Create(imp).Error
}

func (r *AdminRepository) GetUserImport(id uint) (*domain.UserImport, error) {
	var imp domain.UserImport
	if err := r.db.First(&imp, id).Error; err != nil {
		return nil, err
	}
	return &imp, nil
}

func (r *AdminRepository) UpdateUserImport(imp *domain.UserImport) error {
	return r.db.Save(imp).Error
}

// --- Exports ---

func (r *AdminRepository) CreateDataExport(exp *domain.DataExport) error {
	return r.db.Create(exp).Error
}

func (r *AdminRepository) GetDataExport(id uint) (*domain.DataExport, error) {
	var exp domain.DataExport
	if err := r.db.First(&exp, id).Error; err != nil {
		return nil, err
	}
	return &exp, nil
}

func (r *AdminRepository) UpdateDataExport(exp *domain.DataExport) error {
	return r.db.Save(exp).Error
}
