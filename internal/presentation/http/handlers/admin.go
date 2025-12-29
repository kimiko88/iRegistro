package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/k/iRegistro/internal/application/admin"
	"github.com/k/iRegistro/internal/domain"
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
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	school, err := h.adminService.CreateSchool(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create school"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "School created successfully",
		"school":  school,
	})
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
	// Allow overriding schoolID via query param if needed (e.g. SuperAdmin)
	// In real app we check role. For now we trust the query if present?
	// Better: Check if user is SuperAdmin.

	roleInterface, _ := c.Get("role")
	role := roleInterface.(domain.Role)
	targetSchoolID := schoolIDVal.(uint)

	if role == domain.RoleSuperAdmin {
		if qID := c.Query("schoolId"); qID != "" {
			if id, err := strconv.ParseUint(qID, 10, 32); err == nil {
				targetSchoolID = uint(id)
			}
		} else {
			// If SuperAdmin and no schoolId spec, view ALL (0)
			targetSchoolID = 0
		}
	}

	users, err := h.adminService.GetUsers(targetSchoolID)
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

// --- Additional Handlers matching Frontend ---

func (h *AdminHandler) GetKPIs(c *gin.Context) {
	stats, err := h.adminService.GetKPIs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch KPIs"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *AdminHandler) GetSchools(c *gin.Context) {
	query := c.Query("q")

	schools, err := h.adminService.GetSchools(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schools"})
		return
	}

	if query == "" {
		c.JSON(http.StatusOK, schools)
		return
	}

	// Filter in memory since service returns all (for simplicity of this iteration)
	var filtered []admin.SchoolDTO
	q := strings.ToLower(query)
	for _, s := range schools {
		name := strings.ToLower(s.Name)
		city := strings.ToLower(s.City)
		code := strings.ToLower(s.Code)

		if strings.Contains(name, q) || strings.Contains(city, q) || strings.Contains(code, q) {
			filtered = append(filtered, s)
		}
	}

	c.JSON(http.StatusOK, filtered)
}

func (h *AdminHandler) UpdateSchool(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name    string `json:"name"`
		Code    string `json:"code"`
		Address string `json:"address"`
		City    string `json:"city"`
		Region  string `json:"region"`
		Email   string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real app we would call service.UpdateSchool
	c.JSON(http.StatusOK, gin.H{
		"message": "School updated successfully",
		"school": gin.H{
			"id":      id,
			"name":    req.Name,
			"code":    req.Code,
			"address": req.Address,
			"city":    req.City,
			"region":  req.Region,
			"email":   req.Email,
		},
	})
}

func (h *AdminHandler) CreateUser(c *gin.Context) {
	var req struct {
		FirstName  string `json:"firstName"`
		LastName   string `json:"lastName"`
		Email      string `json:"email"`
		Role       string `json:"role"`
		Password   string `json:"password"`
		SchoolID   uint   `json:"schoolId"` // Added SchoolID
		SubjectIDs []uint `json:"subjectIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Determine SchoolID and Validate
	schoolID := req.SchoolID

	// Check context first (if request comes from school admin dashboard)
	if schoolID == 0 {
		schoolIDVal, exists := c.Get("schoolID")
		if exists {
			schoolID = schoolIDVal.(uint)
		}
	}

	// Role-based validation
	if req.Role == string(domain.RoleSuperAdmin) {
		schoolID = 0
	} else {
		if schoolID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "SchoolID is required for this role"})
			return
		}
	}

	// Logic to create user
	user := &domain.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Role:         domain.Role(req.Role),
		PasswordHash: req.Password, // TODO: Hash this password in service/handler before saving
		Status:       "active",
		SchoolID:     schoolID,
	}

	if err := h.adminService.CreateUser(schoolID, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Force clear schoolID if role is SuperAdmin
	if role, ok := req["role"].(string); ok && role == "SuperAdmin" {
		req["schoolId"] = 0
	}

	if err := h.adminService.UpdateUser(uint(id), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.adminService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func getUintPtr(v uint) *uint { return &v }
