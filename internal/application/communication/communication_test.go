package communication

import (
	"testing"
	"time"

	"github.com/k/iRegistro/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// --- Mocks ---

type MockCommRepo struct {
	mock.Mock
}

// Notifications
func (m *MockCommRepo) CreateNotification(n *domain.Notification) error {
	args := m.Called(n)
	return args.Error(0)
}
func (m *MockCommRepo) GetNotificationsByUserID(userID uint, archived bool) ([]domain.Notification, error) {
	args := m.Called(userID, archived)
	return args.Get(0).([]domain.Notification), args.Error(1)
}
func (m *MockCommRepo) MarkNotificationRead(id uint) error {
	return m.Called(id).Error(0)
}
func (m *MockCommRepo) ArchiveNotification(id uint) error {
	return m.Called(id).Error(0)
}
func (m *MockCommRepo) GetPreferences(userID uint) ([]domain.NotificationPreference, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.NotificationPreference), args.Error(1)
}
func (m *MockCommRepo) SavePreferences(prefs []domain.NotificationPreference) error {
	return m.Called(prefs).Error(0)
}

// Messaging
func (m *MockCommRepo) CreateConversation(c *domain.Conversation) error {
	args := m.Called(c)
	// Simulate ID generation
	c.ID = 1
	return args.Error(0)
}
func (m *MockCommRepo) GetConversationsByUserID(userID uint) ([]domain.Conversation, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.Conversation), args.Error(1)
}
func (m *MockCommRepo) GetConversationByID(id uint) (*domain.Conversation, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Conversation), args.Error(1)
}
func (m *MockCommRepo) CreateMessage(msg *domain.Message) error {
	args := m.Called(msg)
	msg.ID = 1
	return args.Error(0)
}
func (m *MockCommRepo) GetMessagesByConversationID(convID uint, limit, offset int) ([]domain.Message, error) {
	args := m.Called(convID, limit, offset)
	return args.Get(0).([]domain.Message), args.Error(1)
}
func (m *MockCommRepo) SoftDeleteMessage(id uint) error {
	return m.Called(id).Error(0)
}

// Colloquiums
func (m *MockCommRepo) CreateColloquiumSlot(slot *domain.ColloquiumSlot) error {
	return m.Called(slot).Error(0)
}
func (m *MockCommRepo) GetAvailableSlots(teacherID uint, from, to time.Time) ([]domain.ColloquiumSlot, error) {
	args := m.Called(teacherID, from, to)
	return args.Get(0).([]domain.ColloquiumSlot), args.Error(1)
}
func (m *MockCommRepo) GetSlotsByTeacherID(teacherID uint) ([]domain.ColloquiumSlot, error) {
	args := m.Called(teacherID)
	return args.Get(0).([]domain.ColloquiumSlot), args.Error(1)
}
func (m *MockCommRepo) GetSlotByID(id uint) (*domain.ColloquiumSlot, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ColloquiumSlot), args.Error(1)
}
func (m *MockCommRepo) UpdateSlot(slot *domain.ColloquiumSlot) error {
	return m.Called(slot).Error(0)
}
func (m *MockCommRepo) DeleteSlot(id uint) error {
	return m.Called(id).Error(0)
}
func (m *MockCommRepo) CreateBooking(booking *domain.ColloquiumBooking) error {
	return m.Called(booking).Error(0)
}
func (m *MockCommRepo) GetBookingsBySlotID(slotID uint) ([]domain.ColloquiumBooking, error) {
	args := m.Called(slotID)
	return args.Get(0).([]domain.ColloquiumBooking), args.Error(1)
}
func (m *MockCommRepo) GetBookingsByParentID(parentID uint) ([]domain.ColloquiumBooking, error) {
	args := m.Called(parentID)
	return args.Get(0).([]domain.ColloquiumBooking), args.Error(1)
}
func (m *MockCommRepo) GetBookingsByDateRange(from, to time.Time) ([]domain.ColloquiumBooking, error) {
	args := m.Called(from, to)
	return args.Get(0).([]domain.ColloquiumBooking), args.Error(1)
}
func (m *MockCommRepo) UpdateBooking(booking *domain.ColloquiumBooking) error {
	return m.Called(booking).Error(0)
}
func (m *MockCommRepo) DeleteBooking(id uint) error {
	return m.Called(id).Error(0)
}

// --- Tests ---

func TestScheduler(t *testing.T) {
	mockRepo := new(MockCommRepo)
	// mock notif service calls via repo mock?
	// Scheduler calls s.notifService.TriggerNotification -> s.repo.GetPreferences + s.repo.CreateNotification
	// This is integration-ish.
	// We'll mock the repo calls that NotificationService makes.

	// Creating real service with mock repo
	notifSvc := NewNotificationService(mockRepo)
	// But we need a Logger
	logger := zap.NewNop()

	sched := NewScheduler(mockRepo, notifSvc, logger)

	// Expectations
	// 1. GetBookingsByDateRange return 1 booking
	mockRepo.On("GetBookingsByDateRange", mock.Anything, mock.Anything).Return([]domain.ColloquiumBooking{{ParentID: 10}}, nil)

	// 2. TriggerNotification flow
	// 2a. GetPreferences (called by TriggerNotification)
	mockRepo.On("GetPreferences", uint(10)).Return([]domain.NotificationPreference{}, nil)
	// 2b. CreateNotification
	mockRepo.On("CreateNotification", mock.MatchedBy(func(n *domain.Notification) bool {
		return n.UserID == 10 && n.Type == domain.NotifTypeSystem && n.Title == "Reminder"
	})).Return(nil)

	// Act
	sched.SendReminders()

	// Assert
	mockRepo.AssertExpectations(t)
}

func TestTriggerNotification(t *testing.T) {
	mockRepo := new(MockCommRepo)
	svc := NewNotificationService(mockRepo)

	// User has Preference for EMAIL (but defaults to InApp logic in service if only EMAIL provided but we force logic check)
	// Service implementation: "if p.Type == notifType { channel = p.Channels[0] }"
	prefs := []domain.NotificationPreference{
		{UserID: 1, Type: domain.NotifTypeGrade, Channels: domain.JSONStringArray{"EMAIL"}},
	}

	mockRepo.On("GetPreferences", uint(1)).Return(prefs, nil)
	mockRepo.On("CreateNotification", mock.MatchedBy(func(n *domain.Notification) bool {
		return n.UserID == 1 && n.Type == domain.NotifTypeGrade && n.Channel == domain.ChannelEmail
	})).Return(nil)

	err := svc.TriggerNotification(1, domain.NotifTypeGrade, "New Grade", "You got an A", nil)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMessagingFlow(t *testing.T) {
	mockRepo := new(MockCommRepo)
	svc := NewMessagingService(mockRepo)

	// Create Conversation
	mockRepo.On("CreateConversation", mock.MatchedBy(func(c *domain.Conversation) bool {
		return c.Subject == "Hello" && len(c.ParticipantIDs) == 2
	})).Return(nil)

	id, err := svc.CreateConversation(1, []uint{2}, "Hello", false)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), id)

	// Send Message
	mockRepo.On("CreateMessage", mock.MatchedBy(func(m *domain.Message) bool {
		return m.ConversationID == 1 && m.Body == "Hi there"
	})).Return(nil)

	msg, err := svc.SendMessage(1, 1, "Hi there", nil)
	assert.NoError(t, err)
	assert.NotNil(t, msg)
}

func TestBookColloquiumSlot(t *testing.T) {
	mockRepo := new(MockCommRepo)
	notifSvc := NewNotificationService(mockRepo) // Not strictly used unless we mock its internal calls, but ColloquiumService uses it.
	// Actually ColloquiumService calls s.notifService... which is a struct.
	// If we want to mock NotifService calls, we'd need an interface for NotifService or just let it run (it uses mockRepo anyway).
	// Since NotifService uses repo, and ColloquiumService uses repo...

	svc := NewColloquiumService(mockRepo, notifSvc)

	slot := &domain.ColloquiumSlot{
		ID: 10, TeacherID: 5, MaxParticipants: 1, IsAvailable: true,
	}

	// Case 1: Success
	mockRepo.On("GetSlotByID", uint(10)).Return(slot, nil)
	// Empty bookings
	mockRepo.On("GetBookingsBySlotID", uint(10)).Return([]domain.ColloquiumBooking{}, nil).Once()
	mockRepo.On("CreateBooking", mock.Anything).Return(nil)

	err := svc.BookSlot(10, 100, 200, "Notes")
	assert.NoError(t, err)

	// Case 2: Full
	mockRepo.On("GetSlotByID", uint(10)).Return(slot, nil)
	// Existing booking
	mockRepo.On("GetBookingsBySlotID", uint(10)).Return([]domain.ColloquiumBooking{{ID: 99}}, nil)

	errFull := svc.BookSlot(10, 101, 201, "Late")
	assert.Error(t, errFull)
	assert.Equal(t, "slot is full", errFull.Error())
}
