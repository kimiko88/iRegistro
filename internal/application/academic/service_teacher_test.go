package academic

import (
	"testing"

	"github.com/k/iRegistro/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockRepo struct {
	mock.Mock
	domain.AcademicRepository // Embed interface
}

func (m *MockRepo) CreateMark(mark *domain.Mark) error {
	args := m.Called(mark)
	return args.Error(0)
}

// Implement other methods to satisfy interface if needed, or stick to what is called
// Since Go interfaces are implicit, we only need to implement what is called IF we pass it as the interface.
// But we need to pass it to NewAcademicService which expects domain.AcademicRepository.
// So we must implement all methods or embed a struct that implements them (stub).
// A better approach for this single test is creating a struct that implements the interface.

type StubRepo struct {
	domain.AcademicRepository
	CreateMarkFn func(mark *domain.Mark) error
}

func (s *StubRepo) CreateMark(mark *domain.Mark) error {
	if s.CreateMarkFn != nil {
		return s.CreateMarkFn(mark)
	}
	return nil
}

type MockNotifier struct {
	domain.NotificationService
	NotifyMarkAddedFn func(mark *domain.Mark)
}

func (m *MockNotifier) NotifyMarkAdded(mark *domain.Mark) {
	if m.NotifyMarkAddedFn != nil {
		m.NotifyMarkAddedFn(mark)
	}
}

func TestCreateMark(t *testing.T) {
	mockRepo := &StubRepo{}
	mockNotifier := &MockNotifier{}
	service := NewAcademicService(mockRepo, nil, mockNotifier)

	t.Run("CreateMark Success and Notification", func(t *testing.T) {
		mark := &domain.Mark{Value: 8, StudentID: 1}

		repoCalled := false
		mockRepo.CreateMarkFn = func(m *domain.Mark) error {
			repoCalled = true
			return nil
		}

		notifCalled := false
		mockNotifier.NotifyMarkAddedFn = func(m *domain.Mark) {
			notifCalled = true
		}

		err := service.CreateMark(mark)
		assert.NoError(t, err)
		assert.True(t, repoCalled)
		assert.True(t, notifCalled)
	})
}
