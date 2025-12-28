package domain

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleSuperAdmin Role = "SuperAdmin"
	RoleAdmin      Role = "Admin"
	RolePrincipal  Role = "Dirigente"
	RoleTeacher    Role = "Insegnante"
	RoleParent     Role = "Genitore"
	RoleStudent    Role = "Studente"
	RoleSecretary  Role = "Segreteria"
)

type User struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	Email          string     `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash   string     `gorm:"not null" json:"-"`
	SchoolID       uint       `gorm:"index" json:"school_id"` // 0 for SuperAdmin
	Role           Role       `gorm:"type:varchar(50);not null" json:"role"`
	FirstName      string     `gorm:"size:100" json:"first_name"`
	LastName       string     `gorm:"size:100" json:"last_name"`
	TwoFAEnabled   bool       `gorm:"default:false" json:"two_fa_enabled"`
	TwoFASecret    string     `gorm:"size:100" json:"-"`
	FailedLogins   int        `gorm:"default:0" json:"-"`
	LockedUntil    *time.Time `json:"locked_until,omitempty"`
	ResetTokenHash string     `gorm:"size:255" json:"-"`
	ResetTokenExp  *time.Time `json:"-"`
	// SPID/CIE Authentication
	AuthMethod      string         `gorm:"type:varchar(20);default:'email'" json:"auth_method"` // email, spid, cie
	SPIDProvider    *string        `gorm:"size:50" json:"spid_provider,omitempty"`
	CIESerialNumber *string        `gorm:"size:50" json:"cie_serial_number,omitempty"`
	ExternalID      *string        `gorm:"size:255;uniqueIndex" json:"-"`
	LastAuthAt      *time.Time     `json:"last_auth_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
