package integration

import (
	"context"
	"testing"
	"time"

	"github.com/k/iRegistro/internal/domain"
	"github.com/k/iRegistro/internal/infrastructure/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	gormPG "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupTestDB(t *testing.T) (*gorm.DB, func()) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db, err := gorm.Open(gormPG.Open(connStr), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	// Auto Migrate
	err = db.AutoMigrate(
		&domain.School{},
		&domain.User{},
		&domain.Campus{},
		&domain.Curriculum{},
		&domain.Class{},
		&domain.Student{},
		&domain.ClassEnrollment{},
	)
	if err != nil {
		t.Fatal(err)
	}

	return db, func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}
}

func TestAcademicRepository_CreateAndGet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, teardown := SetupTestDB(t)
	defer teardown()

	repo := persistence.NewAcademicRepository(db)

	// Seed School
	school := domain.School{Name: "Test School", Code: "TEST01", City: "Test City", Region: "Test Region"}
	db.Create(&school)

	t.Run("Create and Get Class", func(t *testing.T) {
		class := domain.Class{
			Grade:   1,
			Section: "A",
			Year:    "2024/2025",
		}
		err := repo.CreateClass(&class)
		assert.NoError(t, err)
		assert.NotZero(t, class.ID)

		fetched, err := repo.GetClassByID(class.ID)
		assert.NoError(t, err)
		assert.Equal(t, "A", fetched.Section)
	})
}
