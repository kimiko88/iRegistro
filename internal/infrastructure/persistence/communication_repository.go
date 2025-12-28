package persistence

import (
	"time"

	"github.com/k/iRegistro/internal/domain"
	"gorm.io/gorm"
)

type CommunicationRepository struct {
	db *gorm.DB
}

func NewCommunicationRepository(db *gorm.DB) *CommunicationRepository {
	return &CommunicationRepository{db: db}
}

// --- Notifications ---

func (r *CommunicationRepository) CreateNotification(n *domain.Notification) error {
	return r.db.Create(n).Error
}

func (r *CommunicationRepository) GetNotificationsByUserID(userID uint, archived bool) ([]domain.Notification, error) {
	var notifs []domain.Notification
	err := r.db.Where("user_id = ? AND is_archived = ?", userID, archived).
		Order("created_at desc").
		Find(&notifs).Error
	return notifs, err
}

func (r *CommunicationRepository) MarkNotificationRead(id uint) error {
	return r.db.Model(&domain.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

func (r *CommunicationRepository) ArchiveNotification(id uint) error {
	return r.db.Model(&domain.Notification{}).Where("id = ?", id).Update("is_archived", true).Error
}

func (r *CommunicationRepository) GetPreferences(userID uint) ([]domain.NotificationPreference, error) {
	var prefs []domain.NotificationPreference
	err := r.db.Where("user_id = ?", userID).Find(&prefs).Error
	return prefs, err
}

func (r *CommunicationRepository) SavePreferences(prefs []domain.NotificationPreference) error {
	return r.db.Save(prefs).Error
}

// --- Messaging ---

func (r *CommunicationRepository) CreateConversation(c *domain.Conversation) error {
	return r.db.Create(c).Error
}

func (r *CommunicationRepository) GetConversationsByUserID(userID uint) ([]domain.Conversation, error) {
	var convs []domain.Conversation
	// Since ParticipantIDs is JSONB, we need a way to check if userID is in the array.
	// Postgres: @> operator.
	// We need to cast userID to the format stored in JSON. Since it is JSONUintArray, likely [1, 2].
	// However, querying JSONB arrays for a simple value depends on the structure.
	// Assuming ParticipantIDs is stored as `[1, 2]`.
	// Use Gorm datatypes or raw SQL.
	// Raw: "participant_ids @> '[<userID>]'"

	// Construct JSON string for query
	// Note: string concatenation is risky but strictly typed uint is safe.
	// err := r.db.Where("participant_ids @> ?", fmt.Sprintf("[%d]", userID)).Order("last_message_at desc").Find(&convs).Error

	// Safer with Gorm args?
	err := r.db.Where("participant_ids @> ?::jsonb", []uint{userID}).Order("last_message_at desc").Find(&convs).Error
	return convs, err
}

func (r *CommunicationRepository) GetConversationByID(id uint) (*domain.Conversation, error) {
	var c domain.Conversation
	if err := r.db.Preload("Messages").First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CommunicationRepository) CreateMessage(m *domain.Message) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		// Update conversation LastMessageAt
		return tx.Model(&domain.Conversation{}).Where("id = ?", m.ConversationID).
			Update("last_message_at", m.CreatedAt).Error
	})
}

func (r *CommunicationRepository) GetMessagesByConversationID(convID uint, limit, offset int) ([]domain.Message, error) {
	var msgs []domain.Message
	err := r.db.Where("conversation_id = ?", convID).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&msgs).Error
	return msgs, err
}

func (r *CommunicationRepository) SoftDeleteMessage(id uint) error {
	return r.db.Model(&domain.Message{}).Where("id = ?", id).Update("is_deleted", true).Error
}

// --- Colloquiums ---

func (r *CommunicationRepository) CreateColloquiumSlot(slot *domain.ColloquiumSlot) error {
	return r.db.Create(slot).Error
}

func (r *CommunicationRepository) GetAvailableSlots(teacherID uint, from, to time.Time) ([]domain.ColloquiumSlot, error) {
	var slots []domain.ColloquiumSlot
	err := r.db.Where("teacher_id = ? AND date >= ? AND date <= ? AND is_available = ?", teacherID, from, to, true).
		// Optionally check if fully booked?
		// If MaxParticipants > 1, we might need to check count(bookings).
		// Or assume 'IsAvailable' is managed by service.
		Find(&slots).Error
	return slots, err
}

func (r *CommunicationRepository) GetSlotsByTeacherID(teacherID uint) ([]domain.ColloquiumSlot, error) {
	var slots []domain.ColloquiumSlot
	err := r.db.Where("teacher_id = ?", teacherID).Order("date desc, start_time asc").Find(&slots).Error
	return slots, err
}

func (r *CommunicationRepository) GetSlotByID(id uint) (*domain.ColloquiumSlot, error) {
	var slot domain.ColloquiumSlot
	if err := r.db.First(&slot, id).Error; err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *CommunicationRepository) UpdateSlot(slot *domain.ColloquiumSlot) error {
	return r.db.Save(slot).Error
}

func (r *CommunicationRepository) DeleteSlot(id uint) error {
	return r.db.Delete(&domain.ColloquiumSlot{}, id).Error
}

func (r *CommunicationRepository) CreateBooking(booking *domain.ColloquiumBooking) error {
	// Transaction to check availability? Service layer might handle it, but DB constraints help.
	return r.db.Create(booking).Error
}

func (r *CommunicationRepository) GetBookingsBySlotID(slotID uint) ([]domain.ColloquiumBooking, error) {
	var bookings []domain.ColloquiumBooking
	err := r.db.Where("slot_id = ?", slotID).Find(&bookings).Error
	return bookings, err
}

func (r *CommunicationRepository) GetBookingsByParentID(parentID uint) ([]domain.ColloquiumBooking, error) {
	var bookings []domain.ColloquiumBooking
	err := r.db.Preload("Slot").Where("parent_id = ?", parentID).Order("booked_at desc").Find(&bookings).Error
	return bookings, err
}

func (r *CommunicationRepository) UpdateBooking(booking *domain.ColloquiumBooking) error {
	return r.db.Save(booking).Error
}

func (r *CommunicationRepository) DeleteBooking(id uint) error {
	return r.db.Delete(&domain.ColloquiumBooking{}, id).Error
}
