package gdpr

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// AuditLogger handles GDPR access logging
type AuditLogger struct {
	repo AuditRepository
}

type AuditRepository interface {
	LogAccess(log *DataAccessLog) error
	GetAccessLogsForUser(userID uuid.UUID, limit int) ([]DataAccessLog, error)
}

type DataAccessLog struct {
	ID               uuid.UUID
	AccessedByUserID uuid.UUID
	AccessedUserID   *uuid.UUID
	ResourceType     string
	ResourceID       *uuid.UUID
	AccessedAt       time.Time
	IPAddress        string
	UserAgent        string
	Purpose          string
	Action           string
}

func NewAuditLogger(repo AuditRepository) *AuditLogger {
	return &AuditLogger{repo: repo}
}

// LogDataAccess logs access to user data (GDPR Art. 30)
func (a *AuditLogger) LogDataAccess(ctx context.Context, accessedByID, accessedUserID uuid.UUID, resourceType string, resourceID *uuid.UUID, action, purpose string) error {
	log := &DataAccessLog{
		ID:               uuid.New(),
		AccessedByUserID: accessedByID,
		AccessedUserID:   &accessedUserID,
		ResourceType:     resourceType,
		ResourceID:       resourceID,
		AccessedAt:       time.Now(),
		Action:           action,
		Purpose:          purpose,
	}

	// Extract IP and User-Agent from context if available
	if ip, ok := ctx.Value("ip_address").(string); ok {
		log.IPAddress = ip
	}
	if ua, ok := ctx.Value("user_agent").(string); ok {
		log.UserAgent = ua
	}

	return a.repo.LogAccess(log)
}

// GetMyAccessLogs retrieves who accessed a user's data
func (a *AuditLogger) GetMyAccessLogs(userID uuid.UUID, limit int) ([]DataAccessLog, error) {
	return a.repo.GetAccessLogsForUser(userID, limit)
}

// ResourceTypes
const (
	ResourceTypeMarks       = "MARKS"
	ResourceTypeAbsences    = "ABSENCES"
	ResourceTypeMessages    = "MESSAGES"
	ResourceTypeDocuments   = "DOCUMENTS"
	ResourceTypeProfile     = "PROFILE"
	ResourceTypeColloquiums = "COLLOQUIUMS"
)

// Actions
const (
	ActionRead   = "READ"
	ActionWrite  = "WRITE"
	ActionDelete = "DELETE"
	ActionExport = "EXPORT"
)
