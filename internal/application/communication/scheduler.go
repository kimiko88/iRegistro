package communication

import (
	"time"

	"github.com/k/iRegistro/internal/domain"
	"go.uber.org/zap"
)

type Scheduler struct {
	repo         domain.CommunicationRepository
	notifService *NotificationService
	logger       *zap.Logger
}

func NewScheduler(repo domain.CommunicationRepository, notif *NotificationService, logger *zap.Logger) *Scheduler {
	return &Scheduler{repo: repo, notifService: notif, logger: logger}
}

// Start runs a periodic check (e.g., every hour)
func (s *Scheduler) Start() {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			s.SendReminders()
		}
	}()
}

func (s *Scheduler) SendReminders() {
	// Logic: Find Bookings for TOMORROW
	tomorrow := time.Now().AddDate(0, 0, 1)
	startOfDay := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// We need a repository method to get bookings by date range or iterating slots.
	bookings, err := s.repo.GetBookingsByDateRange(startOfDay, endOfDay)
	if err != nil {
		s.logger.Error("Failed to fetch bookings for reminders", zap.Error(err))
		return
	}

	for _, b := range bookings {
		// Stub trigger
		err := s.notifService.TriggerNotification(b.ParentID, domain.NotifTypeSystem, "Reminder", "You have a colloquium tomorrow", nil)
		if err != nil {
			s.logger.Error("Failed to send reminder", zap.Uint("user_id", b.ParentID), zap.Error(err))
		}
	}

}
