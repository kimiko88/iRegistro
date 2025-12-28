package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/k/iRegistro/internal/application/director"
)

type DirectorHandler struct {
	service *director.DirectorService
}

func NewDirectorHandler(service *director.DirectorService) *DirectorHandler {
	return &DirectorHandler{service: service}
}

func (h *DirectorHandler) GetKPIs(c *gin.Context) {
	// schoolID := c.GetUint("schoolID")
	kpis, err := h.service.GetDashboardKPIs(1) // Mock School ID 1
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kpis)
}

func (h *DirectorHandler) GetDocumentsToSign(c *gin.Context) {
	docs, err := h.service.GetPendingDocuments(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, docs)
}

func (h *DirectorHandler) SignDocument(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		PIN string `json:"pin"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SignDocument(uint(id), req.PIN); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
