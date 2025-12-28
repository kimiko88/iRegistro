package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k/iRegistro/internal/gdpr"
)

type GDPRHandler struct {
	complianceSvc *gdpr.ComplianceService
	auditLogger   *gdpr.AuditLogger
}

func NewGDPRHandler(compliance *gdpr.ComplianceService, audit *gdpr.AuditLogger) *GDPRHandler {
	return &GDPRHandler{
		complianceSvc: compliance,
		auditLogger:   audit,
	}
}

// GetUserData returns all user data (GDPR Art. 15 - Right to access)
// GET /users/:userId/data
func (h *GDPRHandler) GetUserData(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Verify user can only access their own data (or admin)
	// This should use proper authorization middleware

	format := c.Query("format")
	if format == "" {
		format = "JSON"
	}

	// Request data export
	export, err := h.complianceSvc.RequestDataExport(c.Request.Context(), userID, format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"export_id": export.ID,
		"status":    export.Status,
		"message":   "Data export initiated. Check status with GET /data-exports/:exportId",
	})
}

// RequestDataExport creates a data export request (GDPR Art. 20 - Data portability)
// POST /users/:userId/data-export
func (h *GDPRHandler) RequestDataExport(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Format string `json:"format" binding:"required,oneof=JSON CSV XML"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	export, err := h.complianceSvc.RequestDataExport(c.Request.Context(), userID, req.Format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, export)
}

// RequestDataDeletion creates a deletion request (GDPR Art. 17 - Right to erasure)
// POST /users/:userId/data-deletion-request
func (h *GDPRHandler) RequestDataDeletion(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deletionReq, err := h.complianceSvc.RequestDataDeletion(c.Request.Context(), userID, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"request_id":            deletionReq.ID,
		"scheduled_deletion_at": deletionReq.ScheduledDeletionAt,
		"message":               "Deletion request submitted. Data will be deleted after 30-day grace period.",
	})
}

// GrantConsent records user consent (GDPR Art. 7)
// PUT /users/:userId/consent/:consentType
func (h *GDPRHandler) GrantConsent(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	consentType := c.Param("consentType")

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	if err := h.complianceSvc.GrantConsent(c.Request.Context(), userID, consentType, ipAddress, userAgent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Consent recorded",
		"consent_type": consentType,
		"granted_at":   c.GetTime("request_time"),
	})
}

// RevokeConsent allows withdrawal of consent (GDPR Art. 7.3)
// DELETE /users/:userId/consent/:consentType
func (h *GDPRHandler) RevokeConsent(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	consentType := c.Param("consentType")

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	if err := h.complianceSvc.RevokeConsent(c.Request.Context(), userID, consentType, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Consent revoked",
		"consent_type": consentType,
	})
}

// GetConsents retrieves all consents for a user
// GET /users/:userId/consent
func (h *GDPRHandler) GetConsents(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	consents, err := h.complianceSvc.GetUserConsents(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, consents)
}

// GetMyAccessLogs shows who accessed user's data
// GET /audit/my-accesses
func (h *GDPRHandler) GetMyAccessLogs(c *gin.Context) {
	// Get current user ID from auth context
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	limit := 100
	logs, err := h.auditLogger.GetMyAccessLogs(userUUID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accesses": logs,
		"count":    len(logs),
	})
}

// SoftDeleteUser performs soft delete (with 30-day grace period)
// DELETE /users/:userId
func (h *GDPRHandler) SoftDeleteUser(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Create deletion request
	deletionReq, err := h.complianceSvc.RequestDataDeletion(c.Request.Context(), userID, "User requested account deletion")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":               "Account scheduled for deletion",
		"scheduled_deletion_at": deletionReq.ScheduledDeletionAt,
		"recovery_period_days":  30,
	})
}
