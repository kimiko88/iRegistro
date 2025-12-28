package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// --- Enums & Types ---

type ClassGroupType string

const (
	GroupArticulate ClassGroupType = "ARTICULATED"
	GroupLanguage   ClassGroupType = "LANGUAGE"
	GroupReligion   ClassGroupType = "RELIGION"
)

type EnrollmentStatus string

const (
	EnrollmentActive      EnrollmentStatus = "ACTIVE"
	EnrollmentTransferred EnrollmentStatus = "TRANSFERRED"
	EnrollmentGraduated   EnrollmentStatus = "GRADUATED"
)

type MarkType string

const (
	MarkNumeric  MarkType = "NUMERIC"
	MarkJudgment MarkType = "JUDGMENT"
)

type AbsenceType string

const (
	AbsenceFull    AbsenceType = "ABSENT" // Full day or specific hour? Usually maps to hour if per-hour
	AbsenceLate    AbsenceType = "LATE"
	AbsenceExcused AbsenceType = "EXCUSED" // Permesso di uscita / Entrata posticipata
	AbsenceDaD     AbsenceType = "DAD"
)

type WeekDay string

const (
	Monday    WeekDay = "MONDAY"
	Tuesday   WeekDay = "TUESDAY"
	Wednesday WeekDay = "WEDNESDAY"
	Thursday  WeekDay = "THURSDAY"
	Friday    WeekDay = "FRIDAY"
	Saturday  WeekDay = "SATURDAY"
	Sunday    WeekDay = "SUNDAY"
)

// --- JSON Types ---

type ScheduleData struct {
	Items []ScheduleItem `json:"items"`
}
type ScheduleItem struct {
	Day       WeekDay `json:"day"`
	Hour      int     `json:"hour"`
	SubjectID uint    `json:"subject_id"`
	TeacherID uint    `json:"teacher_id"`
	Room      string  `json:"room"`
}

// Value/Scan for Gorm JSON support
func (s ScheduleData) Value() (driver.Value, error) {
	return json.Marshal(s)
}
func (s *ScheduleData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &s)
}

// --- Entities ---

type School struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Name           string `gorm:"size:255;not null" json:"name"`
	City           string `gorm:"size:100" json:"city"`
	Region         string `gorm:"size:100" json:"region"`
	VatID          string `gorm:"size:50" json:"vat_id"`
	Address        string `gorm:"size:255" json:"address"`
	PrincipalEmail string `gorm:"size:100" json:"principal_email"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Campus struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	SchoolID  uint   `gorm:"index;not null" json:"school_id"`
	Name      string `gorm:"size:255;not null" json:"name"`
	Address   string `gorm:"size:255" json:"address"`
	CreatedAt time.Time
}

type Curriculum struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	CampusID  uint   `gorm:"index;not null" json:"campus_id"`
	Name      string `gorm:"size:255;not null" json:"name"` // e.g. "Liceo Scientifico"
	Code      string `gorm:"size:50" json:"code"`
	CreatedAt time.Time
}

type Class struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	CurriculumID uint   `gorm:"index;not null" json:"curriculum_id"`
	Grade        int    `json:"grade"`                  // 1, 2, 3, 4, 5
	Section      string `gorm:"size:10" json:"section"` // A, B, C
	Year         string `gorm:"size:20" json:"year"`    // 2024-25
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ClassGroup struct {
	ID      uint           `gorm:"primaryKey" json:"id"`
	ClassID uint           `gorm:"index;not null" json:"class_id"`
	Name    string         `gorm:"size:100" json:"name"` // e.g. "Gruppo Inglese Avanzato"
	Type    ClassGroupType `gorm:"type:varchar(50)" json:"type"`
}

type Subject struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	SchoolID     uint   `gorm:"index;not null" json:"school_id"` // Subjects are defined per school typically
	Code         string `gorm:"size:50;index" json:"code"`       // National Code e.g. MAT01
	Name         string `gorm:"size:255;not null" json:"name"`
	HoursPerWeek int    `json:"hours_per_week"`
}

type ClassSubjectAssignment struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	ClassID   uint       `gorm:"index;not null" json:"class_id"`
	SubjectID uint       `gorm:"index;not null" json:"subject_id"`
	TeacherID uint       `gorm:"index;not null" json:"teacher_id"` // Connects to User (Role=Teacher)
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}

type Student struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SchoolID     uint      `gorm:"index;not null" json:"school_id"`
	FirstName    string    `gorm:"size:100;not null" json:"first_name"`
	LastName     string    `gorm:"size:100;not null" json:"last_name"`
	DateOfBirth  time.Time `json:"date_of_birth"`
	PlaceOfBirth string    `gorm:"size:100" json:"place_of_birth"`
	TaxCode      string    `gorm:"size:16;uniqueIndex" json:"tax_code"`
	Gender       string    `gorm:"size:10" json:"gender"`
	Citizenship  string    `gorm:"size:100" json:"citizenship"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ClassEnrollment struct {
	ID             uint             `gorm:"primaryKey" json:"id"`
	StudentID      uint             `gorm:"index;not null" json:"student_id"`
	ClassID        uint             `gorm:"index;not null" json:"class_id"`
	Year           string           `gorm:"size:20" json:"year"`
	Status         EnrollmentStatus `gorm:"type:varchar(50)" json:"status"`
	EnrollmentDate time.Time        `json:"enrollment_date"`
}

type Mark struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	StudentID     uint           `gorm:"index;not null" json:"student_id"`
	SubjectID     uint           `gorm:"index;not null" json:"subject_id"`
	ClassID       uint           `gorm:"index;not null" json:"class_id"` // Denormalized for easier queries
	TeacherID     uint           `gorm:"index;not null" json:"teacher_id"`
	Value         float64        `gorm:"type:decimal(4,2)" json:"value"` // 1.00 - 10.00
	Type          MarkType       `gorm:"type:varchar(50)" json:"type"`
	Date          time.Time      `json:"date"`
	IsJustified   bool           `gorm:"default:false" json:"is_justified"`
	Justification string         `gorm:"size:255" json:"justification"`
	Weight        float64        `gorm:"default:1.0" json:"weight"` // 100%, 50% etc.
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	DeletedReason string         `gorm:"size:255" json:"deleted_reason"`
}

type Absence struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	StudentID     uint        `gorm:"index;not null" json:"student_id"`
	ClassID       uint        `gorm:"index;not null" json:"class_id"`
	Date          time.Time   `json:"date"` // Day of absence
	Hour          int         `json:"hour"` // 0 for full day, 1-8 for specific hours
	Type          AbsenceType `gorm:"type:varchar(50)" json:"type"`
	IsJustified   bool        `gorm:"default:false" json:"is_justified"`
	JustifiedDate *time.Time  `json:"justified_date"`
	Note          string      `gorm:"size:255" json:"note"`
	CreatedAt     time.Time
}

type Schedule struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	ClassID   uint         `gorm:"uniqueIndex:idx_class_ver" json:"class_id"`
	Version   int          `gorm:"uniqueIndex:idx_class_ver" json:"version"` // To keep history of schedule changes
	Data      ScheduleData `gorm:"type:jsonb" json:"data"`
	ValidFrom time.Time    `json:"valid_from"`
	ValidTo   *time.Time   `json:"valid_to"`
	CreatedAt time.Time
}

type ClassCoordinator struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	TeacherID    uint   `gorm:"index;not null" json:"teacher_id"`
	ClassID      uint   `gorm:"index;not null" json:"class_id"`
	AcademicYear string `gorm:"size:20" json:"academic_year"`
}
