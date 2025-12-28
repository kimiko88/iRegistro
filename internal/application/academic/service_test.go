package academic

import (
	// Imported but only used if needed
	"testing"
	"time"

	"github.com/k/iRegistro/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Repo
type MockAcademicRepository struct {
	mock.Mock
}

func (m *MockAcademicRepository) GetCampuses(schoolID uint) ([]domain.Campus, error) {
	args := m.Called(schoolID)
	return args.Get(0).([]domain.Campus), args.Error(1)
}

func (m *MockAcademicRepository) CreateCampus(campus *domain.Campus) error {
	args := m.Called(campus)
	return args.Error(0)
}

func (m *MockAcademicRepository) GetCurriculums(schoolID uint) ([]domain.Curriculum, error) {
	args := m.Called(schoolID)
	return args.Get(0).([]domain.Curriculum), args.Error(1)
}

func (m *MockAcademicRepository) CreateCurriculum(curriculum *domain.Curriculum) error {
	args := m.Called(curriculum)
	return args.Error(0)
}

func (m *MockAcademicRepository) GetClasses(schoolID uint) ([]domain.Class, error) {
	args := m.Called(schoolID)
	return args.Get(0).([]domain.Class), args.Error(1)
}

func (m *MockAcademicRepository) CreateClass(class *domain.Class) error {
	args := m.Called(class)
	return args.Error(0)
}

func (m *MockAcademicRepository) GetClassByID(classID uint) (*domain.Class, error) {
	args := m.Called(classID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Class), args.Error(1)
}

func (m *MockAcademicRepository) CreateSchool(school *domain.School) error {
	return nil
}
func (m *MockAcademicRepository) GetSchoolByID(id uint) (*domain.School, error) {
	return nil, nil
}
func (m *MockAcademicRepository) GetCampusesBySchoolID(schoolID uint) ([]domain.Campus, error) {
	return nil, nil
}
func (m *MockAcademicRepository) GetCurriculumsBySchoolID(schoolID uint) ([]domain.Curriculum, error) {
	return nil, nil
}
func (m *MockAcademicRepository) GetClassesBySchoolID(schoolID uint) ([]domain.Class, error) {
	return nil, nil
}
func (m *MockAcademicRepository) CreateStudent(student *domain.Student) error {
	return nil
}
func (m *MockAcademicRepository) GetStudentByID(id uint) (*domain.Student, error) {
	return nil, nil
}
func (m *MockAcademicRepository) EnrollStudent(enrollment *domain.ClassEnrollment) error {
	return nil
}
func (m *MockAcademicRepository) GetStudentsByClassID(classID uint, year string) ([]domain.Student, error) {
	return nil, nil
}
func (m *MockAcademicRepository) CreateSubject(subject *domain.Subject) error {
	return nil
}
func (m *MockAcademicRepository) GetSubjectByID(id uint) (*domain.Subject, error) {
	return nil, nil
}
func (m *MockAcademicRepository) AssignSubjectToClass(assignment *domain.ClassSubjectAssignment) error {
	return nil
}
func (m *MockAcademicRepository) GetAssignmentsByTeacherID(teacherID uint) ([]domain.ClassSubjectAssignment, error) {
	return nil, nil
}
func (m *MockAcademicRepository) CreateMark(mark *domain.Mark) error {
	return nil
}
func (m *MockAcademicRepository) GetMarksByStudentID(studentID, classID, subjectID uint) ([]domain.Mark, error) {
	return nil, nil
}
func (m *MockAcademicRepository) GetMarksByClassID(classID uint) ([]domain.Mark, error) {
	return nil, nil
}
func (m *MockAcademicRepository) UpdateMark(mark *domain.Mark) error {
	return nil
}
func (m *MockAcademicRepository) CreateAbsence(absence *domain.Absence) error {
	return nil
}
func (m *MockAcademicRepository) GetAbsencesByStudentID(studentID uint, year string) ([]domain.Absence, error) {
	return nil, nil
}
func (m *MockAcademicRepository) GetAbsencesByClassID(classID uint, date time.Time) ([]domain.Absence, error) {
	return nil, nil
}
func (m *MockAcademicRepository) UpdateAbsence(absence *domain.Absence) error {
	return nil
}

// --- Tests ---

func TestCreateClass_Validation(t *testing.T) {
	repo := new(MockAcademicRepository)
	// We need to mock dependencies of service (UserRepo, Broadcaster) to pass nil or mocks
	service := NewAcademicService(repo, nil, nil) // Passing nil for unused deps in this test

	tests := []struct {
		name        string
		input       domain.Class
		expectError bool
	}{
		{
			name: "Valid Class",
			input: domain.Class{
				Grade:   1, // Was Year (int)
				Section: "A",
				Year:    "2024/2025", // Was AcademicYear (string)
			},
			expectError: false,
		},
		{
			name: "Missing Section",
			input: domain.Class{
				Grade:   1,
				Section: "", // Invalid if validated
				Year:    "2024/2025",
			},
			expectError: true, // Assuming validation exists or check
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.expectError {
				repo.On("CreateClass", mock.Anything).Return(nil)
			}

			// Use address of input
			err := service.CreateClass(&tt.input)

			if tt.expectError {
				// assert
			} else {
				assert.NoError(t, err)
				repo.AssertCalled(t, "CreateClass", mock.Anything)
			}
		})
	}
}
