package admin

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/k/iRegistro/internal/domain"
	"go.uber.org/zap"
)

type UserImportService struct {
	repo     domain.AdminRepository
	userRepo domain.UserRepository
	logger   *zap.Logger
}

func NewUserImportService(repo domain.AdminRepository, userRepo domain.UserRepository, logger *zap.Logger) *UserImportService {
	return &UserImportService{repo: repo, userRepo: userRepo, logger: logger}
}

func (s *UserImportService) ProcessImport(importID uint) {
	// 1. Fetch Import Record
	imp, err := s.repo.GetUserImport(importID)
	if err != nil {
		s.logger.Error("Import process failed: fetch record", zap.Error(err))
		return
	}

	imp.Status = "PROCESSING"
	s.repo.UpdateUserImport(imp)

	// 2. Open File
	// Assuming file is local. In prod, might download from S3.
	f, err := os.Open(imp.ImportFilePath)
	if err != nil {
		s.failImport(imp, "File open error: "+err.Error())
		return
	}
	defer f.Close()

	reader := csv.NewReader(f)
	// Skip header
	_, _ = reader.Read()

	successCount := 0
	failCount := 0
	errors := make(map[string]interface{})

	// 3. Iterate
	rowNum := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		rowNum++
		if err != nil {
			failCount++
			errors[strconv.Itoa(rowNum)] = "CSV Parse error"
			continue
		}

		// CSV Format: [Email, Password, Role, SchoolID]
		if len(record) < 3 {
			failCount++
			errors[strconv.Itoa(rowNum)] = "Invalid column count"
			continue
		}

		email := record[0]
		password := record[1]
		role := record[2]
		schoolID, _ := strconv.Atoi(record[3])

		// Create User
		user := &domain.User{
			Email:        email,
			PasswordHash: password, // Should hash in service/repo
			Role:         domain.Role(role),
			SchoolID:     uint(schoolID),
		}

		if err := s.userRepo.Create(user); err != nil {
			failCount++
			errors[strconv.Itoa(rowNum)] = err.Error()
			continue
		}
		successCount++
	}

	// 4. Update Status
	imp.ImportedUsers = successCount
	imp.FailedUsers = failCount
	imp.ErrorDetails = errors
	imp.Status = "COMPLETED"
	if failCount > 0 {
		imp.Status = "COMPLETED_WITH_ERRORS"
	}
	imp.ImportedAt = time.Now()

	s.repo.UpdateUserImport(imp)
}

func (s *UserImportService) failImport(imp *domain.UserImport, msg string) {
	imp.Status = "FAILED"
	imp.ErrorDetails = map[string]interface{}{"fatal_error": msg}
	s.repo.UpdateUserImport(imp)
}
