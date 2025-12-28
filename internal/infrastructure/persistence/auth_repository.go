package persistence

import (
	"errors"
	"time"

	"github.com/k/iRegistro/internal/domain"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateSession(session *domain.Session) error {
	return r.db.Create(session).Error
}

func (r *AuthRepository) StoreRefreshToken(token *domain.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *AuthRepository) RevokeRefreshToken(tokenHash string) error {
	return r.db.Model(&domain.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked_at", time.Now()).Error
}

func (r *AuthRepository) GetRefreshToken(tokenHash string) (*domain.RefreshToken, error) {
	var token domain.RefreshToken
	if err := r.db.Where("token_hash = ?", tokenHash).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}
