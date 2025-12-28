package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k/iRegistro/internal/application/communication"
	"github.com/k/iRegistro/internal/domain"
)

type CommunicationHandler struct {
	notifService *communication.NotificationService
	msgService   *communication.MessagingService
	colService   *communication.ColloquiumService
}

func NewCommunicationHandler(n *communication.NotificationService, m *communication.MessagingService, c *communication.ColloquiumService) *CommunicationHandler {
	return &CommunicationHandler{notifService: n, msgService: m, colService: c}
}

// --- Notifications ---

func (h *CommunicationHandler) GetNotifications(c *gin.Context) {
	// Get UserID from context (middleware)
	userIDVal, _ := c.Get("userID")
	userID := userIDVal.(uint)

	archived := c.Query("archived") == "true"

	notifs, err := h.notifService.GetUserNotifications(userID, archived)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()}) // Or error status
		return
	}
	c.JSON(http.StatusOK, notifs)
}

func (h *CommunicationHandler) ReadNotification(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.notifService.ReadNotification(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// --- Messaging ---

func (h *CommunicationHandler) CreateConversation(c *gin.Context) {
	var req struct {
		ParticipantIDs []uint `json:"participant_ids"`
		Subject        string `json:"subject"`
		IsGroup        bool   `json:"is_group"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDVal, _ := c.Get("userID")
	initiatorID := userIDVal.(uint)

	id, err := h.msgService.CreateConversation(initiatorID, req.ParticipantIDs, req.Subject, req.IsGroup)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *CommunicationHandler) GetConversations(c *gin.Context) {
	userIDVal, _ := c.Get("userID")
	convs, err := h.msgService.GetUserConversations(userIDVal.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, convs)
}

func (h *CommunicationHandler) SendMessage(c *gin.Context) {
	convID, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Body        string `json:"body"`
		Attachments []struct {
			FilePath string `json:"file_path"`
			FileName string `json:"file_name"`
			FileType string `json:"file_type"`
		} `json:"attachments"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDVal, _ := c.Get("userID")

	// Convert attachments
	var atts []domain.MessageAttachment
	for _, a := range req.Attachments {
		atts = append(atts, domain.MessageAttachment{
			FilePath: a.FilePath,
			FileName: a.FileName,
			FileType: a.FileType,
		})
	}

	msg, err := h.msgService.SendMessage(userIDVal.(uint), uint(convID), req.Body, atts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, msg)
}

func (h *CommunicationHandler) GetMessages(c *gin.Context) {
	convID, _ := strconv.Atoi(c.Param("id"))
	userIDVal, _ := c.Get("userID") // unused in basic impl, but good for validation

	msgs, err := h.msgService.GetConversationMessages(uint(convID), userIDVal.(uint), 50, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, msgs)
}

// --- Colloquiums ---

func (h *CommunicationHandler) CreateSlot(c *gin.Context) {
	var req struct {
		Date            time.Time `json:"date"`
		StartTime       string    `json:"start_time"`
		EndTime         string    `json:"end_time"`
		MaxParticipants int       `json:"max_participants"`
		Type            string    `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDVal, _ := c.Get("userID") // Teacher ID

	err := h.colService.CreateSlot(userIDVal.(uint), req.Date, req.StartTime, req.EndTime, req.MaxParticipants, domain.ColloquiumType(req.Type))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *CommunicationHandler) GetAvailableSlots(c *gin.Context) {
	teacherID, _ := strconv.Atoi(c.Query("teacher_id"))
	slots, err := h.colService.GetAvailableSlots(uint(teacherID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, slots)
}

func (h *CommunicationHandler) BookSlot(c *gin.Context) {
	var req struct {
		SlotID    uint   `json:"slot_id"`
		StudentID uint   `json:"student_id"`
		Notes     string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDVal, _ := c.Get("userID") // Parent ID

	err := h.colService.BookSlot(req.SlotID, userIDVal.(uint), req.StudentID, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}
