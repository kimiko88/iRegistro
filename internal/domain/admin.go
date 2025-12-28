package domain

import (
	"time"
)

// --- Audit Logs ---

type AuditLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SchoolID     *uint     `gorm:"index" json:"school_id,omitempty"` // Nullable for global logs
	UserID       uint      `gorm:"index;not null" json:"user_id"`
	Action       string    `gorm:"size:255;not null" json:"action"`
	ResourceType string    `gorm:"size:100" json:"resource_type"`
	ResourceID   string    `gorm:"size:100" json:"resource_id"`
	Changes      JSONMap   `gorm:"type:jsonb" json:"changes"`
	IPAddress    string    `gorm:"size:50" json:"ip_address"`
	Timestamp    time.Time `json:"timestamp"`
}

// --- School Settings ---

type SchoolSettings struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	SchoolID uint    `gorm:"index;not null;uniqueIndex:idx_school_key" json:"school_id"`
	Key      string  `gorm:"size:100;not null;uniqueIndex:idx_school_key" json:"key"`
	Value    JSONMap `gorm:"type:jsonb" json:"value"`
}

// --- Imports ---

type UserImport struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	SchoolID       uint      `gorm:"index;not null" json:"school_id"`
	ImportFilePath string    `gorm:"size:255" json:"file_path"`
	TotalUsers     int       `json:"total_users"`
	ImportedUsers  int       `json:"imported_users"`
	FailedUsers    int       `json:"failed_users"`
	ErrorDetails   JSONMap   `gorm:"type:jsonb" json:"error_details"`
	Status         string    `gorm:"size:50" json:"status"` // PENDING, PROCESSING, COMPLETED, FAILED
	ImportedAt     time.Time `json:"imported_at"`
}

// --- Exports ---

type DataExport struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SchoolID    uint      `gorm:"index;not null" json:"school_id"`
	ExportType  string    `gorm:"size:10" json:"export_type"` // CSV, JSON
	RequestedBy uint      `gorm:"index" json:"requested_by"`
	FilePath    string    `gorm:"size:255" json:"file_path"`
	Status      string    `gorm:"size:50" json:"status"` // PENDING, READY, FAILED
	CreatedAt   time.Time `json:"created_at"`
	ExpiryDate  time.Time `json:"expiry_date"`
}

// --- Interfaces ---

type AdminRepository interface {
	// Audit
	CreateAuditLog(log *AuditLog) error
	GetAuditLogs(schoolID *uint, limit, offset int) ([]AuditLog, error)

	// Settings
	GetSchoolSettings(schoolID uint) ([]SchoolSettings, error)
	UpsertSchoolSetting(setting *SchoolSettings) error

	// Imports
	CreateUserImport(imp *UserImport) error
	GetUserImport(id uint) (*UserImport, error)
	UpdateUserImport(imp *UserImport) error

	// Exports
	CreateDataExport(exp *DataExport) error
	GetDataExport(id uint) (*DataExport, error)
	UpdateDataExport(exp *DataExport) error
}
