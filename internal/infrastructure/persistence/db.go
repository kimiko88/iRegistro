package persistence

import (
	"fmt"

	"github.com/k/iRegistro/internal/config"
	"github.com/k/iRegistro/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Drop views and problematic tables that block migration
	db.Exec("DROP VIEW IF EXISTS teacher_workload; DROP VIEW IF EXISTS pdp_documents; DROP VIEW IF EXISTS student_current_marks; DROP VIEW IF EXISTS class_statistics; DROP VIEW IF EXISTS colloquium_availability; DROP VIEW IF EXISTS pdp_documents CASCADE;")
	db.Exec("DROP TABLE IF EXISTS sessions CASCADE;")

	// Auto-migrate
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Session{},
		&domain.RefreshToken{},
		&domain.School{},
		&domain.Campus{},
		&domain.Curriculum{},
		&domain.Class{},
		&domain.ClassGroup{},
		&domain.Subject{},
		&domain.ClassSubjectAssignment{},
		&domain.Student{},
		&domain.ClassEnrollment{},
		&domain.Mark{},
		&domain.Absence{},
		&domain.Schedule{},
		&domain.ClassCoordinator{},
	)
	if err != nil {
		return nil, fmt.Errorf("auto-migrate failed: %w", err)
	}

	// Basic Seeding
	seedInitialData(db)

	return db, nil
}

func seedInitialData(db *gorm.DB) {
	// Seed School
	var count int64
	db.Model(&domain.School{}).Count(&count)
	if count == 0 {
		fmt.Println("Seeding initial school...")
		db.Create(&domain.School{ID: 1, Name: "Test School", Code: "TEST001", City: "Roma"})
	}

	// Seed Campus
	db.Model(&domain.Campus{}).Count(&count)
	if count == 0 {
		fmt.Println("Seeding initial campus...")
		db.Create(&domain.Campus{ID: 1, SchoolID: 1, Name: "Main Campus"})
	}

	// Seed Curriculum
	db.Model(&domain.Curriculum{}).Count(&count)
	if count == 0 {
		fmt.Println("Seeding initial curriculum...")
		db.Create(&domain.Curriculum{ID: 1, SchoolID: 1, Name: "Standard Curriculum", Code: "STD01"})
	}

	// Ensure specific test users are updated...
	fmt.Println("Updating test user roles...")
	db.Model(&domain.User{}).Where("email = ?", "admin@test.it").Update("role", "Admin")
	db.Model(&domain.User{}).Where("email = ?", "secretary@test.it").Update("role", "Secretary")
	db.Model(&domain.User{}).Where("email = ?", "teacher@test.it").Update("role", "Teacher")
}
