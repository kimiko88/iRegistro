package director

import (
	"testing"
	"time"

	"github.com/k/iRegistro/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepo implements both AcademicRepository and ReportingRepository partially
type MockRepo struct {
	mock.Mock
}

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

// Stub other methods needed for compilation if interface check is strict?
// No, Mock object can just be passed if we don't assign it to a full interface variable in test setup,
// OR we rely on partial interface satisfaction if we define small interfaces in service.
// But Service uses `domain.AcademicRepository`. We need full stubs or helper struct.
// Let's use a struct that embeds mock.Mock and implements minimal required,
// but Go requires all methods.
// Easier: Import `secretary.MockRepo` (if public) or duplicate stubs?
// Duplicating minimal stubs to satisfy compiler.

func (m *MockRepo) CreateSchool(school *domain.School) error                     { return nil }
func (m *MockRepo) GetSchoolByID(id uint) (*domain.School, error)                { return nil, nil }
func (m *MockRepo) CreateCampus(campus *domain.Campus) error                     { return nil }
func (m *MockRepo) GetCampusesBySchoolID(schoolID uint) ([]domain.Campus, error) { return nil, nil }
func (m *MockRepo) CreateCurriculum(curriculum *domain.Curriculum) error         { return nil }
func (m *MockRepo) GetCurriculumsBySchoolID(schoolID uint) ([]domain.Curriculum, error) {
	return nil, nil
}
func (m *MockRepo) CreateClass(class *domain.Class) error                      { return nil }
func (m *MockRepo) GetClassByID(id uint) (*domain.Class, error)                { return nil, nil }
func (m *MockRepo) GetClassesBySchoolID(schoolID uint) ([]domain.Class, error) { return nil, nil }
func (m *MockRepo) CreateStudent(student *domain.Student) error                { return nil }
func (m *MockRepo) GetStudentByID(id uint) (*domain.Student, error)            { return nil, nil }
func (m *MockRepo) EnrollStudent(enrollment *domain.ClassEnrollment) error     { return nil }
func (m *MockRepo) GetStudentsByClassID(classID uint, year string) ([]domain.Student, error) {
	return nil, nil
}
func (m *MockRepo) CreateSubject(subject *domain.Subject) error                          { return nil }
func (m *MockRepo) GetSubjectByID(id uint) (*domain.Subject, error)                      { return nil, nil }
func (m *MockRepo) AssignSubjectToClass(assignment *domain.ClassSubjectAssignment) error { return nil }
func (m *MockRepo) GetAssignmentsByTeacherID(teacherID uint) ([]domain.ClassSubjectAssignment, error) {
	return nil, nil
}
func (m *MockRepo) CreateMark(mark *domain.Mark) error { return nil }
func (m *MockRepo) GetMarksByStudentID(studentID uint, classID uint, subjectID uint) ([]domain.Mark, error) {
	return nil, nil
}
func (m *MockRepo) GetMarksByClassID(classID uint) ([]domain.Mark, error) { return nil, nil }
func (m *MockRepo) UpdateMark(mark *domain.Mark) error                    { return nil }
func (m *MockRepo) CreateAbsence(absence *domain.Absence) error           { return nil }
func (m *MockRepo) GetAbsencesByStudentID(studentID uint, year string) ([]domain.Absence, error) {
	return nil, nil
}
func (m *MockRepo) GetAbsencesByClassID(classID uint, date time.Time) ([]domain.Absence, error) {
	return nil, nil
}
func (m *MockRepo) UpdateAbsence(absence *domain.Absence) error { return nil }

// Reporting methods stubs
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
func (m *MockRepo) DeleteDocument(id uint) error { return nil }
func (m *MockRepo) CountDocumentsByStatus(schoolID uint, status domain.DocumentStatus) (int64, error) {
	return 0, nil
}
func (m *MockRepo) CountDocumentsUpdatedSince(schoolID uint, status []domain.DocumentStatus, since time.Time) (int64, error) {
	return 0, nil
}

// Actual Tests

func TestGetDashboardKPIs(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewDirectorService(mockRepo, mockRepo)

	kpis, err := svc.GetDashboardKPIs(1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1200), kpis.TotalStudents)
}

func TestSignDocument(t *testing.T) {
	mockRepo := new(MockRepo)
	svc := NewDirectorService(mockRepo, mockRepo)

	doc := &domain.Document{ID: 10, Status: domain.DocStatusDraft}

	mockRepo.On("GetDocumentByID", uint(10)).Return(doc, nil)
	mockRepo.On("AddSignature", mock.Anything).Return(nil)
	mockRepo.On("UpdateDocument", mock.MatchedBy(func(d *domain.Document) bool {
		return d.Status == domain.DocStatusSigned
	})).Return(nil)

	err := svc.SignDocument(10, "123456")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
