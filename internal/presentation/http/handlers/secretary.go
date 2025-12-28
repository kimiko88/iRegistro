package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/k/iRegistro/internal/application/secretary"
)

type SecretaryHandler struct {
	service *secretary.SecretaryService
}

func NewSecretaryHandler(service *secretary.SecretaryService) *SecretaryHandler {
	return &SecretaryHandler{service: service}
}

func (h *SecretaryHandler) GetInbox(c *gin.Context) {
	// SchoolID from context (middleware should set it or User's school)
	// Assuming User struct has SchoolID
	// For now hardcoded or retrieved from user claims
	schoolID := uint(1) // Placeholder

	docs, err := h.service.GetInbox(schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, docs)
}

func (h *SecretaryHandler) GetArchive(c *gin.Context) {
	schoolID := uint(1) // Placeholder
	docs, err := h.service.GetArchive(schoolID, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, docs)
}

func (h *SecretaryHandler) ApproveDocument(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := c.GetUint("userID")

	if err := h.service.ApproveDocument(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *SecretaryHandler) RejectDocument(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Reason string `json:"reason"`
	}
	c.BindJSON(&req)

	if err := h.service.RejectDocument(uint(id), req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *SecretaryHandler) BatchPrint(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pdfBytes, err := h.service.BatchPrint(req.IDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=batch.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func (h *SecretaryHandler) GetDashboardStats(c *gin.Context) {
	schoolID := uint(1) // Placeholder, replace with actual context user school ID
	stats, err := h.service.GetDashboardStats(schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
