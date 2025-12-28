package admin

import (
	"time"

	"github.com/k/iRegistro/internal/domain"
)

type AuditService struct {
	repo domain.AdminRepository
}

func NewAuditService(repo domain.AdminRepository) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) LogAction(schoolID *uint, userID uint, action, resType, resID, ip string, changes domain.JSONMap) error {
	log := &domain.AuditLog{
		SchoolID:     schoolID,
		UserID:       userID,
		Action:       action,
		ResourceType: resType,
		ResourceID:   resID,
		IPAddress:    ip,
		Changes:      changes,
		Timestamp:    time.Now(),
	}
	return s.repo.CreateAuditLog(log)
}

func (s *AuditService) GetLogs(schoolID *uint, limit, offset int) ([]domain.AuditLog, error) {
	return s.repo.GetAuditLogs(schoolID, limit, offset)
}
