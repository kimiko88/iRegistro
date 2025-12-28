package communication

import (
	"time"

	"github.com/k/iRegistro/internal/domain"
)

type MessagingService struct {
	repo domain.CommunicationRepository
	// notifService *NotificationService // To notify on new message
}

func NewMessagingService(repo domain.CommunicationRepository) *MessagingService {
	return &MessagingService{repo: repo}
}

func (s *MessagingService) CreateConversation(initiatorID uint, participantIDs []uint, subject string, isGroup bool) (uint, error) {
	// Validate participants?
	// Add initiator to participants if not present
	found := false
	for _, p := range participantIDs {
		if p == initiatorID {
			found = true
			break
		}
	}
	if !found {
		participantIDs = append(participantIDs, initiatorID)
	}

	convType := domain.ConvTypeOneToOne
	if isGroup {
		convType = domain.ConvTypeClass // Simplification for now
	}

	conv := &domain.Conversation{
		Type:           convType,
		Subject:        subject,
		ParticipantIDs: domain.JSONUintArray(participantIDs),
		CreatedAt:      time.Now(),
		LastMessageAt:  time.Now(),
	}

	if err := s.repo.CreateConversation(conv); err != nil {
		return 0, err
	}

	return conv.ID, nil
}

func (s *MessagingService) SendMessage(senderID, convID uint, body string, attachments []domain.MessageAttachment) (*domain.Message, error) {
	// Verify sender is participant? (Requires GetConversationByID + logic)

	msg := &domain.Message{
		ConversationID: convID,
		SenderID:       senderID,
		Body:           body,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
	}

	// Convert attachments to JSONMapArray (map string interface)
	if len(attachments) > 0 {
		var atts domain.JSONMapArray
		for _, a := range attachments {
			atts = append(atts, map[string]interface{}{
				"file_path": a.FilePath,
				"file_name": a.FileName,
				"file_type": a.FileType,
			})
		}
		msg.Attachments = atts
	}

	if err := s.repo.CreateMessage(msg); err != nil {
		return nil, err
	}

	// TODO: Trigger Notification to other participants
	return msg, nil
}

func (s *MessagingService) GetUserConversations(userID uint) ([]domain.Conversation, error) {
	return s.repo.GetConversationsByUserID(userID)
}

func (s *MessagingService) GetConversationMessages(convID uint, userID uint, limit, offset int) ([]domain.Message, error) {
	// Verify participation?
	return s.repo.GetMessagesByConversationID(convID, limit, offset)
}

func (s *MessagingService) SoftDeleteMessage(msgID uint, userID uint) error {
	// Verify ownership?
	return s.repo.SoftDeleteMessage(msgID)
}
