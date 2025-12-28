package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// --- Enums ---

type DocumentType string

const (
	DocReportCard  DocumentType = "PAGELLA"
	Doc15May       DocumentType = "DOCUMENTO_15_MAGGIO"
	DocPDP         DocumentType = "PDP"
	DocPFI         DocumentType = "PFI"
	DocPFP         DocumentType = "PFP"
	DocPCTO        DocumentType = "PCTO"
	DocOrientation DocumentType = "ORIENTAMENTO"
)

type DocumentStatus string

const (
	DocStatusDraft    DocumentStatus = "DRAFT"
	DocStatusSigned   DocumentStatus = "SIGNED"
	DocStatusArchived DocumentStatus = "ARCHIVED"
)

type PCTOStatus string

const (
	PCTOActive    PCTOStatus = "ACTIVE"
	PCTOCompleted PCTOStatus = "COMPLETED"
	PCTOSuspended PCTOStatus = "SUSPENDED"
)

// --- JSON Types ---

// Generic JSON map for flexible document data
type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}
func (j *JSONMap) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &j)
}

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}
func (a *StringArray) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

// --- Entities ---

type Document struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	SchoolID     uint           `gorm:"index;not null" json:"school_id"`
	Type         DocumentType   `gorm:"type:varchar(50);index" json:"type"`
	Title        string         `gorm:"size:255" json:"title"`
	Data         JSONMap        `gorm:"type:jsonb" json:"data"` // Flexible payload
	CreatedBy    uint           `json:"created_by"`             // UserID
	CreatedAt    time.Time      `json:"created_at"`
	AcademicYear string         `gorm:"size:20" json:"academic_year"`
	Status       DocumentStatus `gorm:"type:varchar(50);default:'DRAFT'" json:"status"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations (Polymorphic-ish, usually linked via Data or separate tables,
	// but for SQL querying it helps to have explicit FKs if common)
	StudentID *uint `gorm:"index" json:"student_id,omitempty"`
	ClassID   *uint `gorm:"index" json:"class_id,omitempty"`

	Signatures []DocumentSignature `json:"signatures"`
}

type DocumentSignature struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	DocumentID         uint      `gorm:"index;not null" json:"document_id"`
	SignerID           uint      `gorm:"index;not null" json:"signer_id"` // UserID
	SignatureTimestamp time.Time `json:"signature_timestamp"`
	IPAddress          string    `gorm:"size:50" json:"ip_address"`
	IsValid            bool      `gorm:"default:true" json:"is_valid"`
}

// --- PCTO ---

type PCTOProject struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	SchoolID    uint    `gorm:"index;not null" json:"school_id"`
	Name        string  `gorm:"size:255;not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Companies   JSONMap `gorm:"type:jsonb" json:"companies"` // Array of company details
	CreatedAt   time.Time
}

type PCTOAssignment struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	PCTOProjectID uint       `gorm:"index;not null" json:"pcto_project_id"`
	StudentID     uint       `gorm:"index;not null" json:"student_id"`
	CompanyID     string     `gorm:"size:100" json:"company_id"` // External ID or name from Companies JSON
	StartDate     time.Time  `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	HoursPlanned  int        `json:"hours_planned"`
	Status        PCTOStatus `gorm:"type:varchar(50)" json:"status"`

	Hours []PCTOHour `json:"hours"`
}

type PCTOHour struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	AssignmentID uint   `gorm:"index;not null" json:"assignment_id"`
	Week         string `gorm:"size:20" json:"week"` // e.g. "2024-W12"
	HoursWorked  int    `json:"hours_worked"`
	Description  string `gorm:"size:255" json:"description"`
	TutorNotes   string `gorm:"type:text" json:"tutor_notes"`
	Status       string `gorm:"size:50" json:"status"` // PENDING, APPROVED
}

// --- Orientation ---

type OrientationActivity struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SchoolID    uint      `gorm:"index;not null" json:"school_id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Date        time.Time `json:"date"`
	Hours       int       `json:"hours"`
	TargetGrade int       `json:"target_grade"` // e.g. 5 for final year
}

type OrientationParticipation struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ActivityID  uint   `gorm:"index;not null" json:"activity_id"`
	StudentID   uint   `gorm:"index;not null" json:"student_id"`
	HoursEarned int    `json:"hours_earned"`
	Evaluation  string `gorm:"size:255" json:"evaluation"`
}

// --- PDP Specifics (Stored in Document.Data usually, but if specific queries needed...) ---
// Keeping PDP data in Document.Data as JSON for flexibility as per request requirements.
