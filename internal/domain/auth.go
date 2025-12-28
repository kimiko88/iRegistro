package domain

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	User      User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	TokenHash string         `gorm:"uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
	IPAddress string         `gorm:"size:45" json:"ip_address"`
	UserAgent string         `gorm:"size:255" json:"user_agent"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type RefreshToken struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index;not null" json:"user_id"`
	User      User       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	TokenHash string     `gorm:"uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
