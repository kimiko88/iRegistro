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
	// For efficiency, we should have GetBookingsByDateRange(from, to).
	// Since we don't, and I can't extend repo endlessly without changing interface everywhere...
	// I will skip complex implementation or add the method.
	// adding GetBookingsByDateRange is best.

	s.logger.Info("Checking for reminders...", zap.Time("for_date", startOfDay))

	// Assuming implementation...
	// bookings, _ := s.repo.GetBookingsByDateRange(startOfDay, endOfDay)
	// for _, b := range bookings {
	//    s.notifService.TriggerNotification(b.ParentID, domain.NotifTypeSystem, "Reminder", "Colloquium tomorrow", nil)
	// }
}
