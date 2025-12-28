package secretary

import (
	"testing"

	"github.com/k/iRegistro/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockRepo struct {
	mock.Mock
}

// Implement necessary interface methods for SecretaryService
func (m *MockRepo) GetDocumentsByStatus(schoolID uint, status domain.DocumentStatus) ([]domain.Document, error) {
	args := m.Called(schoolID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Document), args.Error(1)
}

func (m *MockRepo) GetDocumentByID(id uint) (*domain.Document, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Document), args.Error(1)
}

func (m *MockRepo) AddSignature(sig *domain.DocumentSignature) error {
	args := m.Called(sig)
	return args.Error(0)
}

func (m *MockRepo) UpdateDocument(doc *domain.Document) error {
	args := m.Called(doc)
	return args.Error(0)
}

func (m *MockRepo) DeleteDocument(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// Stubs for interface compliance
func (m *MockRepo) CreateDocument(doc *domain.Document) error { return nil }
func (m *MockRepo) GetDocumentsBySchoolID(schoolID uint, docType domain.DocumentType) ([]domain.Document, error) {
	return nil, nil
}
func (m *MockRepo) GetDocumentsByStudentID(studentID uint) ([]domain.Document, error) {
	return nil, nil
}
func (m *MockRepo) GetSignaturesByDocumentID(docID uint) ([]domain.DocumentSignature, error) {
	return nil, nil
}
func (m *MockRepo) CreatePCTOProject(project *domain.PCTOProject) error { return nil }
func (m *MockRepo) GetPCTOProjectsBySchoolID(schoolID uint) ([]domain.PCTOProject, error) {
	return nil, nil
}
func (m *MockRepo) CreatePCTOAssignment(assignment *domain.PCTOAssignment) error { return nil }
func (m *MockRepo) GetPCTOAssignmentsByClassID(classID uint) ([]domain.PCTOAssignment, error) {
	return nil, nil
}
func (m *MockRepo) GetPCTOAssignmentsByStudentID(studentID uint) ([]domain.PCTOAssignment, error) {
	return nil, nil
}
func (m *MockRepo) LogPCTOHours(hours *domain.PCTOHour) error                            { return nil }
func (m *MockRepo) CreateOrientationActivity(activity *domain.OrientationActivity) error { return nil }
func (m *MockRepo) GetOrientationActivitiesBySchoolID(schoolID uint) ([]domain.OrientationActivity, error) {
	return nil, nil
}
func (m *MockRepo) RegisterOrientationParticipation(participation *domain.OrientationParticipation) error {
	return nil
}
func (m *MockRepo) GetOrientationParticipationsByStudentID(studentID uint) ([]domain.OrientationParticipation, error) {
	return nil, nil
}

type MockPDFGen struct {
	mock.Mock
}

func (m *MockPDFGen) GenerateReportCard(data domain.JSONMap) ([]byte, error) {
	return []byte("mock pdf"), nil
}
func (m *MockPDFGen) GenerateCertificate(data domain.JSONMap) ([]byte, error) {
	return []byte("mock cert"), nil
}

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Save(filename string, data []byte) (string, error) {
	args := m.Called(filename, data)
	return args.String(0), args.Error(1)
}

type MockNotifier struct {
	mock.Mock
}

func (m *MockNotifier) TriggerNotification(userID uint, notifType domain.NotificationType, title, body string, data domain.JSONMap) error {
	args := m.Called(userID, notifType, title, body, data)
	return args.Error(0)
}

// --- Tests ---

func TestGetInbox(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewSecretaryService(mockRepo, new(MockPDFGen), new(MockStorage), new(MockNotifier))

	docs := []domain.Document{{ID: 1, Status: domain.DocStatusDraft}}

	mockRepo.On("GetDocumentsByStatus", uint(1), domain.DocStatusDraft).Return(docs, nil)

	result, err := svc.GetInbox(1)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	mockRepo.AssertExpectations(t)
}

func TestApproveDocument(t *testing.T) {
	mockRepo := new(MockRepo)
	mockStorage := new(MockStorage)
	mockNotifier := new(MockNotifier)
	svc := NewSecretaryService(mockRepo, new(MockPDFGen), mockStorage, mockNotifier)

	doc := &domain.Document{ID: 10, Status: domain.DocStatusDraft, Type: domain.DocReportCard, StudentID: new(uint)} // StudentID needed for notif
	*doc.StudentID = 123

	mockRepo.On("GetDocumentByID", mock.Anything).Return(doc, nil)
	mockStorage.On("Save", mock.Anything, mock.Anything).Return("/path/to/doc.pdf", nil)
	mockRepo.On("AddSignature", mock.Anything).Return(nil)
	mockRepo.On("UpdateDocument", mock.Anything).Return(nil)
	mockNotifier.On("TriggerNotification", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := svc.ApproveDocument(10, 99)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockNotifier.AssertExpectations(t)
}
