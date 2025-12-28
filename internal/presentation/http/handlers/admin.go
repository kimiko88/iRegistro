package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/k/iRegistro/internal/application/admin"
)

type AdminHandler struct {
	adminService  *admin.AdminService
	auditService  *admin.AuditService
	importService *admin.UserImportService
	exportService *admin.DataExportService
}

func NewAdminHandler(adm *admin.AdminService, audit *admin.AuditService, imp *admin.UserImportService, exp *admin.DataExportService) *AdminHandler {
	return &AdminHandler{adminService: adm, auditService: audit, importService: imp, exportService: exp}
}

// --- SuperAdmin ---

func (h *AdminHandler) CreateSchool(c *gin.Context) {
	// Only for SuperAdmin (Middleware should check this)
	var req struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDVal, _ := c.Get("userID")

	// For now, service implementation is stubbed, so just return OK
	c.JSON(http.StatusCreated, gin.H{"message": "School created", "details": req})

	h.auditService.LogAction(nil, userIDVal.(uint), "CREATE_SCHOOL", "SCHOOL", "", c.ClientIP(), nil)
}

// --- School Admin ---

func (h *AdminHandler) GetSettings(c *gin.Context) {
	schoolIDVal, _ := c.Get("schoolID")
	settings, err := h.adminService.GetSchoolSettings(schoolIDVal.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, settings)
}

func (h *AdminHandler) GetUsers(c *gin.Context) {
	schoolIDVal, _ := c.Get("schoolID")
	users, err := h.adminService.GetUsers(schoolIDVal.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *AdminHandler) UpdateSetting(c *gin.Context) {
	var req struct {
		Key   string                 `json:"key"`
		Value map[string]interface{} `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolIDVal, _ := c.Get("schoolID")
	userIDVal, _ := c.Get("userID")

	if err := h.adminService.UpdateSchoolSetting(schoolIDVal.(uint), userIDVal.(uint), req.Key, req.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *AdminHandler) RequestExport(c *gin.Context) {
	var req struct {
		Format string `json:"format"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schoolIDVal, _ := c.Get("schoolID")
	userIDVal, _ := c.Get("userID")

	id, err := h.exportService.RequestExport(schoolIDVal.(uint), userIDVal.(uint), req.Format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"export_id": id})
}

func (h *AdminHandler) GetAuditLogs(c *gin.Context) {
	schoolIDVal, _ := c.Get("schoolID")

	// Query params: limit, offset

	logs, err := h.auditService.GetLogs(getUintPtr(schoolIDVal.(uint)), 50, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

func getUintPtr(v uint) *uint { return &v }
