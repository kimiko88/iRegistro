package persistence

import (
	"context"
	"errors"

	"github.com/k/iRegistro/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.db.Preload("Subjects").Preload("School").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindAll(schoolID uint) ([]domain.User, error) {
	var users []domain.User
	query := r.db
	if schoolID != 0 {
		query = query.Where("school_id = ?", schoolID)
	}
	if err := query.Preload("Subjects").Preload("School").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetByExternalID finds a user by external ID (SPID/CIE identifier)
func (r *UserRepository) GetByExternalID(ctx context.Context, externalID string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("external_id = ?", externalID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	// Use Association to replace subjects
	if err := r.db.Model(user).Association("Subjects").Replace(user.Subjects); err != nil {
		return err
	}
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}

func (r *UserRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Count(&count).Error
	return count, err
}

func (r *UserRepository) CountBySchoolAndRole(schoolID uint, role domain.Role) (int64, error) {
	var count int64
	err := r.db.Model(&domain.User{}).Where("school_id = ? AND role = ?", schoolID, role).Count(&count).Error
	return count, err
}
