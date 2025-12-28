package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=registro password=password dbname=registro port=5433 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Attempting to add campus_id column to curriculums table...")
	err = db.Exec("ALTER TABLE curriculums ADD COLUMN IF NOT EXISTS campus_id bigint;").Error
	if err != nil {
		fmt.Printf("Error adding column: %v\n", err)
	} else {
		fmt.Println("Successfully added campus_id column (or it already existed).")
	}

	fmt.Println("Standardizing all user roles to English...")
	roleMap := map[string]string{
		"Segreteria": "Secretary",
		"Insegnante": "Teacher",
		"Docente":    "Teacher",
		"Dirigente":  "Principal",
		"Genitore":   "Parent",
		"Studente":   "Student",
		"admin":      "Admin", // Ensure case consistency if needed
	}

	for oldRole, newRole := range roleMap {
		err = db.Exec("UPDATE users SET role = ? WHERE role = ?;", newRole, oldRole).Error
		if err != nil {
			fmt.Printf("Error updating role %s: %v\n", oldRole, err)
		} else {
			fmt.Printf("Standardized %s -> %s\n", oldRole, newRole)
		}
	}

	// Double check specific test users
	fmt.Println("Ensuring specific test users are updated...")
	db.Exec("UPDATE users SET role = 'Admin' WHERE email = 'admin@test.it';")
	db.Exec("UPDATE users SET role = 'Secretary' WHERE email = 'secretary@test.it';")
	db.Exec("UPDATE users SET role = 'Teacher' WHERE email = 'teacher1@test.it';")
}
