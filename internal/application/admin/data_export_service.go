package admin

import (
	"encoding/json"
	"os"
	"time"

	"github.com/k/iRegistro/internal/domain"
)

type DataExportService struct {
	repo domain.AdminRepository
}

func NewDataExportService(repo domain.AdminRepository) *DataExportService {
	return &DataExportService{repo: repo}
}

func (s *DataExportService) RequestExport(schoolID, userID uint, format string) (uint, error) {
	exp := &domain.DataExport{
		SchoolID:    schoolID,
		RequestedBy: userID,
		ExportType:  format,
		Status:      "PENDING",
		CreatedAt:   time.Now(),
		ExpiryDate:  time.Now().Add(24 * time.Hour), // 24h retention
	}

	if err := s.repo.CreateDataExport(exp); err != nil {
		return 0, err
	}

	// Trigger generation async
	go s.generateExport(exp)

	return exp.ID, nil
}

func (s *DataExportService) generateExport(exp *domain.DataExport) {
	// Stub generation logic
	// 1. Fetch all data for school (Users, Classes, etc.)
	// 2. Serialize to CSV/JSON
	// 3. Write to temp file

	// Mocking success
	data := map[string]string{"sample": "data"}
	file, _ := os.CreateTemp("", "export_*.json")
	defer file.Close()
	json.NewEncoder(file).Encode(data)

	exp.FilePath = file.Name()
	exp.Status = "READY"
	s.repo.UpdateDataExport(exp)
}
