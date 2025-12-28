package reporting

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

func (m *MockRepo) CreateDocument(doc *domain.Document) error {
	args := m.Called(doc)
	return args.Error(0)
}
func (m *MockRepo) GetDocumentByID(id uint) (*domain.Document, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Document), args.Error(1)
}
func (m *MockRepo) UpdateDocument(doc *domain.Document) error {
	args := m.Called(doc)
	return args.Error(0)
}
func (m *MockRepo) AddSignature(sig *domain.DocumentSignature) error {
	args := m.Called(sig)
	return args.Error(0)
}

// Stubs for others
func (m *MockRepo) GetDocumentsBySchoolID(schoolID uint, docType domain.DocumentType) ([]domain.Document, error) {
	return nil, nil
}
func (m *MockRepo) GetDocumentsByStudentID(studentID uint) ([]domain.Document, error) {
	return nil, nil
}
func (m *MockRepo) DeleteDocument(id uint) error { return nil }
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
	args := m.Called(data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

// --- Tests ---

func TestGetDocumentPDF(t *testing.T) {
	mockRepo := new(MockRepo)
	mockPDF := new(MockPDFGen)
	svc := NewReportingService(mockRepo, mockPDF)

	doc := &domain.Document{
		ID:   1,
		Type: domain.DocReportCard,
		Data: domain.JSONMap{"foo": "bar"},
	}

	// Expectations
	mockRepo.On("GetDocumentByID", uint(1)).Return(doc, nil)
	mockPDF.On("GenerateReportCard", doc.Data).Return([]byte("%PDF..."), nil)

	// Act
	pdf, err := svc.GetDocumentPDF(1)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, []byte("%PDF..."), pdf)
	mockRepo.AssertExpectations(t)
	mockPDF.AssertExpectations(t)
}

func TestSignDocument(t *testing.T) {
	mockRepo := new(MockRepo)
	mockPDF := new(MockPDFGen)
	svc := NewReportingService(mockRepo, mockPDF)

	doc := &domain.Document{
		ID:     1,
		Status: domain.DocStatusDraft,
	}

	// Expectations
	mockRepo.On("GetDocumentByID", uint(1)).Return(doc, nil)
	mockRepo.On("AddSignature", mock.MatchedBy(func(sig *domain.DocumentSignature) bool {
		return sig.DocumentID == 1 && sig.SignerID == 99 && sig.IsValid
	})).Return(nil)

	// Expect status update
	mockRepo.On("UpdateDocument", mock.MatchedBy(func(d *domain.Document) bool {
		return d.ID == 1 && d.Status == domain.DocStatusSigned
	})).Return(nil)

	// Act
	err := svc.SignDocument(1, 99, "127.0.0.1")

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
