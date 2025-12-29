package main

import (
	"fmt"
	"log"

	"github.com/k/iRegistro/internal/config"
	"github.com/k/iRegistro/internal/domain"
	"github.com/k/iRegistro/internal/infrastructure/persistence"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := persistence.NewDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	email := "admin@test.it"
	password := "password123"
	role := domain.RoleSuperAdmin

	var user domain.User
	result := db.Where("email = ?", email).First(&user)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	if result.Error != nil {
		// Create
		fmt.Printf("User %s not found, creating...\n", email)
		user = domain.User{
			Email:        email,
			PasswordHash: string(hash),
			Role:         role,
			FirstName:    "Admin",
			LastName:     "Test",
			SchoolID:     1,
		}
		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
		fmt.Println("User created successfully.")
	} else {
		// Update
		fmt.Printf("User %s found, updating...\n", email)
		user.PasswordHash = string(hash)
		user.Role = role
		user.SchoolID = 1 // Linked to Test School to satisfy FK, but role is SuperAdmin
		if err := db.Save(&user).Error; err != nil {
			log.Fatalf("Failed to update user: %v", err)
		}
		fmt.Println("User updated successfully.")
	}
}
