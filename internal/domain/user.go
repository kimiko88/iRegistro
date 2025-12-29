package domain

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleSuperAdmin Role = "SuperAdmin"
	RoleAdmin      Role = "Admin"
	RolePrincipal  Role = "Principal"
	RoleTeacher    Role = "Teacher"
	RoleParent     Role = "Parent"
	RoleStudent    Role = "Student"
	RoleSecretary  Role = "Secretary"
)

type User struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	Email          string     `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash   string     `gorm:"not null" json:"-"`
	SchoolID       uint       `gorm:"index" json:"schoolId"` // 0 for SuperAdmin
	Role           Role       `gorm:"type:varchar(50);not null" json:"role"`
	Subjects       []Subject  `gorm:"many2many:user_subjects;" json:"subjects,omitempty"`
	Status         string     `gorm:"type:varchar(20);default:'active'" json:"status"` // active, inactive
	FirstName      string     `gorm:"size:100" json:"firstName"`
	LastName       string     `gorm:"size:100" json:"lastName"`
	TwoFAEnabled   bool       `gorm:"default:false" json:"twoFaEnabled"`
	TwoFASecret    string     `gorm:"size:100" json:"-"`
	FailedLogins   int        `gorm:"default:0" json:"-"`
	LockedUntil    *time.Time `json:"lockedUntil,omitempty"`
	ResetTokenHash string     `gorm:"size:255" json:"-"`
	ResetTokenExp  *time.Time `json:"-"`
	// SPID/CIE Authentication
	AuthMethod      string         `gorm:"type:varchar(20);default:'email'" json:"authMethod"` // email, spid, cie
	SPIDProvider    *string        `gorm:"size:50" json:"spidProvider,omitempty"`
	CIESerialNumber *string        `gorm:"size:50" json:"cieSerialNumber,omitempty"`
	ExternalID      *string        `gorm:"size:255;uniqueIndex" json:"-"`
	LastAuthAt      *time.Time     `json:"lastAuthAt,omitempty"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
