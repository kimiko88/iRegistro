package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// --- Notifications ---

type NotificationType string
type NotificationChannel string

const (
	NotifTypeGrade      NotificationType = "GRADE"
	NotifTypeAbsence    NotificationType = "ABSENCE"
	NotifTypeGeneral    NotificationType = "GENERAL"
	NotifTypeColloquium NotificationType = "COLLOQUIUM"
	NotifTypeSystem     NotificationType = "SYSTEM"

	ChannelEmail NotificationChannel = "EMAIL"
	ChannelSMS   NotificationChannel = "SMS"
	ChannelPush  NotificationChannel = "PUSH"
	ChannelInApp NotificationChannel = "IN_APP"
)

type Notification struct {
	ID         uint                `gorm:"primaryKey" json:"id"`
	UserID     uint                `gorm:"index;not null" json:"user_id"`
	Type       NotificationType    `gorm:"size:50;not null" json:"type"`
	Title      string              `gorm:"size:255;not null" json:"title"`
	Body       string              `gorm:"type:text;not null" json:"body"`
	Data       JSONMap             `gorm:"type:jsonb" json:"data"` // e.g., {"grade_id": 123}
	Channel    NotificationChannel `gorm:"size:50" json:"channel"`
	IsRead     bool                `gorm:"default:false" json:"is_read"`
	IsArchived bool                `gorm:"default:false" json:"is_archived"`
	CreatedAt  time.Time           `json:"created_at"`
}

type NotificationPreference struct {
	UserID   uint             `gorm:"primaryKey" json:"user_id"`
	Type     NotificationType `gorm:"primaryKey;size:50" json:"type"`
	Channels JSONStringArray  `gorm:"type:jsonb" json:"channels"` // e.g. ["EMAIL", "PUSH"]
}

// --- Messaging ---

type ConversationType string

const (
	ConvTypeOneToOne ConversationType = "ONE_TO_ONE"
	ConvTypeClass    ConversationType = "CLASS" // Teacher -> All Parents
	ConvTypeStaff    ConversationType = "STAFF"
)

type Conversation struct {
	ID             uint             `gorm:"primaryKey" json:"id"`
	Type           ConversationType `gorm:"size:50;default:'ONE_TO_ONE'" json:"type"`
	Subject        string           `gorm:"size:255" json:"subject"`         // For class comms
	ClassID        *uint            `gorm:"index" json:"class_id,omitempty"` // Optional link to class
	ParticipantIDs JSONUintArray    `gorm:"type:jsonb" json:"participant_ids"`
	LastMessageAt  time.Time        `json:"last_message_at"`
	CreatedAt      time.Time        `json:"created_at"`

	Messages []Message `gorm:"foreignKey:ConversationID" json:"messages,omitempty"`
}

type Message struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	ConversationID uint         `gorm:"index;not null" json:"conversation_id"`
	SenderID       uint         `gorm:"index;not null" json:"sender_id"`
	Body           string       `gorm:"type:text;not null" json:"body"`
	Attachments    JSONMapArray `gorm:"type:jsonb" json:"attachments"` // Array of {file_path, file_name, file_type}
	IsDeleted      bool         `gorm:"default:false" json:"is_deleted"`
	EditedAt       *time.Time   `json:"edited_at,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
}

// MessageAttachment helper struct for JSONB
type MessageAttachment struct {
	FilePath string `json:"file_path"`
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
}

// --- Colloquiums ---

type ColloquiumType string

const (
	ColloquiumIndividual ColloquiumType = "INDIVIDUA"
	ColloquiumGeneral    ColloquiumType = "GENERALE"
)

type ColloquiumSlot struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	TeacherID       uint           `gorm:"index;not null" json:"teacher_id"`
	ClassID         *uint          `gorm:"index" json:"class_id,omitempty"` // Optional hint
	Date            time.Time      `gorm:"type:date;not null" json:"date"`
	StartTime       string         `gorm:"Type:varchar(5);not null" json:"start_time"` // HH:MM
	EndTime         string         `gorm:"Type:varchar(5);not null" json:"end_time"`   // HH:MM
	MaxParticipants int            `gorm:"default:1" json:"max_participants"`
	Type            ColloquiumType `gorm:"default:'INDIVIDUA'" json:"type"`
	IsAvailable     bool           `gorm:"default:true" json:"is_available"`
	CreatedAt       time.Time      `json:"created_at"`

	Bookings []ColloquiumBooking `gorm:"foreignKey:SlotID" json:"bookings,omitempty"`
}

type ColloquiumBooking struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	SlotID         uint      `gorm:"index;not null" json:"slot_id"`
	ParentID       uint      `gorm:"index;not null" json:"parent_id"`
	StudentID      uint      `gorm:"index;not null" json:"student_id"`
	BookedAt       time.Time `json:"booked_at"`
	NotesBefore    string    `gorm:"type:text" json:"notes_before"`
	NotesAfter     string    `gorm:"type:text" json:"notes_after"`
	FeedbackRating *int      `json:"feedback_rating"` // 1-5
	FeedbackText   string    `gorm:"type:text" json:"feedback_text"`

	Slot ColloquiumSlot `gorm:"foreignKey:SlotID" json:"slot,omitempty"`
}

// --- Interfaces ---

type CommunicationRepository interface {
	// Notifications
	CreateNotification(n *Notification) error
	GetNotificationsByUserID(userID uint, archived bool) ([]Notification, error)
	MarkNotificationRead(id uint) error
	ArchiveNotification(id uint) error
	GetPreferences(userID uint) ([]NotificationPreference, error)
	SavePreferences(prefs []NotificationPreference) error

	// Messaging
	CreateConversation(c *Conversation) error
	GetConversationsByUserID(userID uint) ([]Conversation, error)
	GetConversationByID(id uint) (*Conversation, error)
	CreateMessage(m *Message) error
	GetMessagesByConversationID(convID uint, limit, offset int) ([]Message, error)
	SoftDeleteMessage(id uint) error

	// Colloquiums
	CreateColloquiumSlot(slot *ColloquiumSlot) error
	GetAvailableSlots(teacherID uint, from, to time.Time) ([]ColloquiumSlot, error)
	GetSlotsByTeacherID(teacherID uint) ([]ColloquiumSlot, error)
	GetSlotByID(id uint) (*ColloquiumSlot, error)
	UpdateSlot(slot *ColloquiumSlot) error
	DeleteSlot(id uint) error

	CreateBooking(booking *ColloquiumBooking) error
	GetBookingsBySlotID(slotID uint) ([]ColloquiumBooking, error)
	GetBookingsByParentID(parentID uint) ([]ColloquiumBooking, error)
	GetBookingsByDateRange(from, to time.Time) ([]ColloquiumBooking, error)
	UpdateBooking(booking *ColloquiumBooking) error
	DeleteBooking(id uint) error
}

// --- Helpers for GORM JSONB ---

type JSONUintArray []uint

func (a JSONUintArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JSONUintArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, a)
}

type JSONStringArray []string

func (a JSONStringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JSONStringArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, a)
}

type JSONMapArray []map[string]interface{}

func (a JSONMapArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *JSONMapArray) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, a)
}
