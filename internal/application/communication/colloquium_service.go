package communication

import (
	"errors"
	"time"

	"github.com/k/iRegistro/internal/domain"
)

type ColloquiumService struct {
	repo         domain.CommunicationRepository
	notifService *NotificationService
}

func NewColloquiumService(repo domain.CommunicationRepository, notif *NotificationService) *ColloquiumService {
	return &ColloquiumService{repo: repo, notifService: notif}
}

func (s *ColloquiumService) CreateSlot(teacherID uint, date time.Time, start, end string, maxParticipants int, cType domain.ColloquiumType) error {
	// Basic validation
	if maxParticipants < 1 {
		maxParticipants = 1
	}

	slot := &domain.ColloquiumSlot{
		TeacherID:       teacherID,
		Date:            date,
		StartTime:       start,
		EndTime:         end,
		MaxParticipants: maxParticipants,
		Type:            cType,
		IsAvailable:     true,
		CreatedAt:       time.Now(),
	}

	return s.repo.CreateColloquiumSlot(slot)
}

func (s *ColloquiumService) GetAvailableSlots(teacherID uint) ([]domain.ColloquiumSlot, error) {
	// Return future slots
	today := time.Now()
	nextMonth := today.AddDate(0, 1, 0)
	return s.repo.GetAvailableSlots(teacherID, today, nextMonth)
}

func (s *ColloquiumService) BookSlot(slotID, parentID, studentID uint, notes string) error {
	// 1. Get Slot
	slot, err := s.repo.GetSlotByID(slotID)
	if err != nil {
		return err
	}

	// 2. Check if available
	if !slot.IsAvailable {
		return errors.New("slot is not available")
	}

	// 3. Check Capacity
	bookings, err := s.repo.GetBookingsBySlotID(slotID)
	if err != nil {
		return err
	}

	if len(bookings) >= slot.MaxParticipants {
		return errors.New("slot is full")
	}

	booking := &domain.ColloquiumBooking{
		SlotID:      slotID,
		ParentID:    parentID,
		StudentID:   studentID,
		BookedAt:    time.Now(),
		NotesBefore: notes,
	}

	if err := s.repo.CreateBooking(booking); err != nil {
		return err
	}

	// 2. Notify Teacher (Optional: now we have teacherID from slot)
	// s.notifService.TriggerNotification(slot.TeacherID, ...)

	return nil
}

func (s *ColloquiumService) GetParentBookings(parentID uint) ([]domain.ColloquiumBooking, error) {
	return s.repo.GetBookingsByParentID(parentID)
}

func (s *ColloquiumService) AddFeedback(bookingID uint, rating int, text string) error {
	booking := &domain.ColloquiumBooking{
		ID:             bookingID,
		FeedbackRating: &rating,
		FeedbackText:   text,
	}
	return s.repo.UpdateBooking(booking)
}
