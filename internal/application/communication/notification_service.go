package communication

import (
	"time"

	"github.com/k/iRegistro/internal/domain"
)

type NotificationService struct {
	repo domain.CommunicationRepository
}

func NewNotificationService(repo domain.CommunicationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

// TriggerNotification sends a notification based on user preferences.
func (s *NotificationService) TriggerNotification(userID uint, notifType domain.NotificationType, title, body string, data domain.JSONMap) error {
	// 1. Get User Preferences
	prefs, err := s.repo.GetPreferences(userID)
	if err != nil {
		return err
	}

	// 2. Determine Channel (Default to InApp if no pref)
	channel := domain.ChannelInApp
	for _, p := range prefs {
		if p.Type == notifType && len(p.Channels) > 0 {
			// For simplicity, pick the first configured channel.
			// Real implementation might broadcast to all.
			// TODO: Support multiple channels
			channel = domain.NotificationChannel(p.Channels[0])
			break
		}
	}

	// 3. Create Notification Record
	n := &domain.Notification{
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Body:      body,
		Data:      data,
		Channel:   channel,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateNotification(n); err != nil {
		return err
	}

	// 4. Send to External Provider (Email/SMS/Push)
	// TODO: Integrate actual providers.
	go s.sendExternal(n)

	return nil
}

func (s *NotificationService) sendExternal(n *domain.Notification) {
	// Stub for sending email/sms/push
	// logger.Info("Sending notification", "channel", n.Channel, "user", n.UserID)
}

func (s *NotificationService) GetUserNotifications(userID uint, archived bool) ([]domain.Notification, error) {
	return s.repo.GetNotificationsByUserID(userID, archived)
}

func (s *NotificationService) ReadNotification(id uint) error {
	return s.repo.MarkNotificationRead(id)
}

func (s *NotificationService) ArchiveNotification(id uint) error {
	return s.repo.ArchiveNotification(id)
}

func (s *NotificationService) UpdatePreferences(userID uint, prefs []domain.NotificationPreference) error {
	// Validate?
	return s.repo.SavePreferences(prefs)
}
