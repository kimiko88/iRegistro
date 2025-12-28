package gdpr

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// ComplianceService handles GDPR user rights
type ComplianceService struct {
	repo            ComplianceRepository
	auditLogger     *AuditLogger
	encryptionSvc   *EncryptionService
	storageBasePath string
}

type ComplianceRepository interface {
	// Consent management
	SaveConsent(consent *UserConsent) error
	GetConsents(userID uuid.UUID) ([]UserConsent, error)
	RevokeConsent(userID uuid.UUID, consentType string, reason string) error

	// Data export
	CreateDataExport(export *DataExport) error
	GetDataExport(exportID uuid.UUID) (*DataExport, error)
	UpdateDataExport(export *DataExport) error

	// Data deletion
	CreateDeletionRequest(req *DataDeletionRequest) error
	GetDeletionRequest(reqID uuid.UUID) (*DataDeletionRequest, error)
	ApproveDeletionRequest(reqID, approverID uuid.UUID) error
	SoftDeleteUser(userID uuid.UUID) error
	HardDeleteUser(userID uuid.UUID) error

	// Data retrieval for export
	GetUserCompleteData(userID uuid.UUID) (*UserCompleteData, error)
}

type UserConsent struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	ConsentType   string
	GivenAt       time.Time
	Granted       bool
	IPAddress     string
	UserAgent     string
	RevokedAt     *time.Time
	RevokedReason string
}

type DataExport struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	RequestedAt  time.Time
	Format       string // JSON, CSV, XML
	FilePath     string
	FileSize     int64
	ExpiryDate   time.Time
	DownloadedAt *time.Time
	Status       string
	ErrorMessage string
}

type DataDeletionRequest struct {
	ID                  uuid.UUID
	UserID              uuid.UUID
	RequestedAt         time.Time
	ScheduledDeletionAt *time.Time
	Reason              string
	Status              string
	ProcessedByUserID   *uuid.UUID
	CompletedAt         *time.Time
	RejectionReason     string
}

type UserCompleteData struct {
	UserID     uuid.UUID              `json:"user_id" xml:"user_id"`
	Profile    map[string]interface{} `json:"profile" xml:"profile"`
	Marks      []interface{}          `json:"marks" xml:"marks>mark"`
	Absences   []interface{}          `json:"absences" xml:"absences>absence"`
	Messages   []interface{}          `json:"messages" xml:"messages>message"`
	Documents  []interface{}          `json:"documents" xml:"documents>document"`
	ExportedAt time.Time              `json:"exported_at" xml:"exported_at"`
}

// Consent types
const (
	ConsentCommunications = "COMMUNICATIONS"
	ConsentPhotos         = "PHOTOS"
	ConsentBiometric      = "BIOMETRIC"
	ConsentDataProcessing = "DATA_PROCESSING"
)

func NewComplianceService(repo ComplianceRepository, audit *AuditLogger, enc *EncryptionService, storagePath string) *ComplianceService {
	return &ComplianceService{
		repo:            repo,
		auditLogger:     audit,
		encryptionSvc:   enc,
		storageBasePath: storagePath,
	}
}

// GrantConsent records user consent (GDPR Art. 7)
func (s *ComplianceService) GrantConsent(ctx context.Context, userID uuid.UUID, consentType string, ipAddress, userAgent string) error {
	consent := &UserConsent{
		ID:          uuid.New(),
		UserID:      userID,
		ConsentType: consentType,
		GivenAt:     time.Now(),
		Granted:     true,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
	}

	if err := s.repo.SaveConsent(consent); err != nil {
		return err
	}

	// Audit log
	return s.auditLogger.LogDataAccess(ctx, userID, userID, "CONSENT", &consent.ID, ActionWrite, "Consent granted")
}

// RevokeConsent allows user to withdraw consent (GDPR Art. 7.3)
func (s *ComplianceService) RevokeConsent(ctx context.Context, userID uuid.UUID, consentType, reason string) error {
	if err := s.repo.RevokeConsent(userID, consentType, reason); err != nil {
		return err
	}

	return s.auditLogger.LogDataAccess(ctx, userID, userID, "CONSENT", nil, ActionWrite, "Consent revoked")
}

// GetUserConsents retrieves all consents for a user
func (s *ComplianceService) GetUserConsents(userID uuid.UUID) ([]UserConsent, error) {
	return s.repo.GetConsents(userID)
}

// RequestDataExport handles right to data portability (GDPR Art. 20)
func (s *ComplianceService) RequestDataExport(ctx context.Context, userID uuid.UUID, format string) (*DataExport, error) {
	export := &DataExport{
		ID:          uuid.New(),
		UserID:      userID,
		RequestedAt: time.Now(),
		Format:      format,
		Status:      "PENDING",
		ExpiryDate:  time.Now().Add(30 * 24 * time.Hour), // 30 days
	}

	if err := s.repo.CreateDataExport(export); err != nil {
		return nil, err
	}

	// Process export asynchronously in production
	go s.processDataExport(export)

	return export, s.auditLogger.LogDataAccess(ctx, userID, userID, "DATA_EXPORT", &export.ID, ActionExport, "User data export requested")
}

// processDataExport generates the export file
func (s *ComplianceService) processDataExport(export *DataExport) error {
	export.Status = "PROCESSING"
	s.repo.UpdateDataExport(export)

	// Retrieve complete user data
	data, err := s.repo.GetUserCompleteData(export.UserID)
	if err != nil {
		export.Status = "ERROR"
		export.ErrorMessage = err.Error()
		s.repo.UpdateDataExport(export)
		return err
	}

	data.ExportedAt = time.Now()

	// Generate file based on format
	filename := fmt.Sprintf("user_data_%s_%s.%s", export.UserID, time.Now().Format("20060102"), export.Format)
	filePath := filepath.Join(s.storageBasePath, "exports", filename)

	// Ensure directory exists
	os.MkdirAll(filepath.Dir(filePath), 0755)

	var fileErr error
	switch export.Format {
	case "JSON":
		fileErr = s.exportToJSON(data, filePath)
	case "CSV":
		fileErr = s.exportToCSV(data, filePath)
	case "XML":
		fileErr = s.exportToXML(data, filePath)
	default:
		fileErr = errors.New("unsupported format")
	}

	if fileErr != nil {
		export.Status = "ERROR"
		export.ErrorMessage = fileErr.Error()
	} else {
		export.Status = "READY"
		export.FilePath = filePath
		// Get file size
		if info, err := os.Stat(filePath); err == nil {
			export.FileSize = info.Size()
		}
	}

	return s.repo.UpdateDataExport(export)
}

func (s *ComplianceService) exportToJSON(data *UserCompleteData, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (s *ComplianceService) exportToCSV(data *UserCompleteData, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers and data for each section
	// This is a simplified version - in production, you'd create separate CSV files
	writer.Write([]string{"Section", "Data"})
	writer.Write([]string{"Profile", fmt.Sprintf("%+v", data.Profile)})
	for _, mark := range data.Marks {
		writer.Write([]string{"Mark", fmt.Sprintf("%+v", mark)})
	}

	return nil
}

func (s *ComplianceService) exportToXML(data *UserCompleteData, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")
	return encoder.Encode(data)
}

// RequestDataDeletion handles right to erasure (GDPR Art. 17)
func (s *ComplianceService) RequestDataDeletion(ctx context.Context, userID uuid.UUID, reason string) (*DataDeletionRequest, error) {
	req := &DataDeletionRequest{
		ID:                  uuid.New(),
		UserID:              userID,
		RequestedAt:         time.Now(),
		ScheduledDeletionAt: timePtr(time.Now().Add(30 * 24 * time.Hour)), // 30-day grace period
		Reason:              reason,
		Status:              "PENDING",
	}

	if err := s.repo.CreateDeletionRequest(req); err != nil {
		return nil, err
	}

	return req, s.auditLogger.LogDataAccess(ctx, userID, userID, "DELETION_REQUEST", &req.ID, ActionWrite, "Data deletion requested")
}

// ApproveDeletion approves a deletion request (school admin only)
func (s *ComplianceService) ApproveDeletion(ctx context.Context, requestID, approverID uuid.UUID) error {
	req, err := s.repo.GetDeletionRequest(requestID)
	if err != nil {
		return err
	}

	// Soft delete user data
	if err := s.repo.SoftDeleteUser(req.UserID); err != nil {
		return err
	}

	if err := s.repo.ApproveDeletionRequest(requestID, approverID); err != nil {
		return err
	}

	return s.auditLogger.LogDataAccess(ctx, approverID, req.UserID, "DELETION_REQUEST", &requestID, ActionDelete, "Data deletion approved")
}

// HardDeleteExpiredUsers permanently deletes users after grace period
func (s *ComplianceService) HardDeleteExpiredUsers() error {
	// This should be run as a cron job
	// Find all approved deletion requests past their scheduled deletion date
	// and perform hard delete
	return nil
}

func timePtr(t time.Time) *time.Time {
	return &t
}
