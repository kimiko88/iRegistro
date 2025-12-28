package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/k/iRegistro/internal/application/reporting"
	"github.com/k/iRegistro/internal/domain"
)

type ReportingHandler struct {
	service *reporting.ReportingService
}

func NewReportingHandler(service *reporting.ReportingService) *ReportingHandler {
	return &ReportingHandler{service: service}
}

// --- Documents ---

func (h *ReportingHandler) GetDocuments(c *gin.Context) {
	schoolIDStr := c.Param("schoolId")
	schoolID, _ := strconv.Atoi(schoolIDStr)
	docType := c.Query("type")

	docs, err := h.service.GetDocumentsBySchoolID(uint(schoolID), domain.DocumentType(docType))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, docs)
}

func (h *ReportingHandler) GenerateReportCard(c *gin.Context) {
	// Post data: ClassID, AcademicYear, maybe simple JSON payload
	var req struct {
		ClassID      uint           `json:"class_id"`
		AcademicYear string         `json:"academic_year"`
		Data         domain.JSONMap `json:"data"` // Optional: override data? Or service calculates?
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolIDStr := c.Param("schoolId")
	schoolID, _ := strconv.Atoi(schoolIDStr)

	// Create for single student or loop for class?
	// Handler name implies class generation?
	// For MVP creating one generic document or loop inside service.
	// Assume single document creation for now or stub.

	// Assuming logic generates 1 document for simplicity of this handler example
	doc, err := h.service.CreateReportCard(uint(schoolID), 0, req.ClassID, req.AcademicYear, req.Data, 1) // 1=Creator (TODO: Get from context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, doc)
}

func (h *ReportingHandler) SignDocument(c *gin.Context) {
	docIDStr := c.Param("documentId")
	docID, _ := strconv.Atoi(docIDStr)

	// Get UserID from Context (AuthMiddleware)
	userID := 1 // Stub

	if err := h.service.SignDocument(uint(docID), uint(userID), c.ClientIP()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Document signed"})
}

func (h *ReportingHandler) GetDocumentPDF(c *gin.Context) {
	docIDStr := c.Param("documentId")
	docID, _ := strconv.Atoi(docIDStr)

	pdfBytes, err := h.service.GetDocumentPDF(uint(docID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=document_%d.pdf", docID))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// --- PCTO ---

func (h *ReportingHandler) GetPCTOProjects(c *gin.Context) {
	schoolIDStr := c.Param("schoolId")
	schoolID, _ := strconv.Atoi(schoolIDStr)

	projects, err := h.service.GetPCTOProjects(uint(schoolID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projects)
}

func (h *ReportingHandler) CreatePCTOProject(c *gin.Context) {
	var project domain.PCTOProject
	if err := c.BindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolIDStr := c.Param("schoolId")
	schoolID, _ := strconv.Atoi(schoolIDStr)
	project.SchoolID = uint(schoolID)

	if err := h.service.CreatePCTOProject(&project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, project)
}

// --- Orientation ---
func (h *ReportingHandler) RegisterOrientation(c *gin.Context) {
	var part domain.OrientationParticipation
	if err := c.BindJSON(&part); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.RegisterOrientation(&part); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, part)
}
