package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k/iRegistro/internal/application/academic"
	"github.com/k/iRegistro/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock objects for TeacherHandler
type MockAcademicService struct {
	mock.Mock
}

// Implement necessary methods...
// Since AcademicService is a concrete struct in the handler, we might need to interface it or mock dependencies of service.
// The TeacherHandler takes *academic.AcademicService.
// Testing handlers with concrete services that depend on repositories requires mocking repositories.
// Since we don't have an interface for AcademicService injected, we'd need to integration test or refactor.
// Given strict instructions to "Add test", and previous patterns, let's see if we can use existing MockRepo for Academic.

// However, constructing a real AcademicService with Mock Repos is better.

type MockAcademicRepo struct {
	mock.Mock
}

// Stub methods...
func (m *MockAcademicRepo) CreateClass(c *domain.Class) error                    { return nil }
func (m *MockAcademicRepo) GetClassesBySchoolID(id uint) ([]domain.Class, error) { return nil, nil }
func (m *MockAcademicRepo) GetClassByID(id uint) (*domain.Class, error)          { return nil, nil }

// ... (Add others as needed or use interface)
// The interface is domain.AcademicRepository.

// Let's rely on the fact that we can't easily mock the concrete service struct without deep instantiation.
// Or we can assume the user is okay with a basic test that checks routing/parsing if we can't mock service easily?
// No, let's try to verify `CreateMark` logic at least.

// Wait, I can't redefine MockAcademicRepo if it exists in another package's test file.
// I'll check if `internal/application/academic/service_test.go` defines mocks I can import? No, usually in `_test.go` they are package private or same package.
// I'll try to define a minimal Mock here or skip deep logic.

// Actually, `NewTeacherHandler` takes `*academic.AcademicService`.
// So I must instantiate `AcademicService` with a mock repository.

type MockRepoForTeacher struct {
	mock.Mock
}

// School/Campus
func (m *MockRepoForTeacher) CreateSchool(school *domain.School) error      { return nil }
func (m *MockRepoForTeacher) GetSchoolByID(id uint) (*domain.School, error) { return nil, nil }
func (m *MockRepoForTeacher) CreateCampus(campus *domain.Campus) error      { return nil }
func (m *MockRepoForTeacher) GetCampusesBySchoolID(schoolID uint) ([]domain.Campus, error) {
	return nil, nil
}

// Curriculum/Class
func (m *MockRepoForTeacher) CreateCurriculum(curriculum *domain.Curriculum) error { return nil }
func (m *MockRepoForTeacher) GetCurriculumsBySchoolID(schoolID uint) ([]domain.Curriculum, error) {
	return nil, nil
}
func (m *MockRepoForTeacher) CreateClass(class *domain.Class) error       { return nil }
func (m *MockRepoForTeacher) GetClassByID(id uint) (*domain.Class, error) { return nil, nil }
func (m *MockRepoForTeacher) GetClassesBySchoolID(schoolID uint) ([]domain.Class, error) {
	return nil, nil
}

// Student
func (m *MockRepoForTeacher) CreateStudent(student *domain.Student) error            { return nil }
func (m *MockRepoForTeacher) GetStudentByID(id uint) (*domain.Student, error)        { return nil, nil }
func (m *MockRepoForTeacher) EnrollStudent(enrollment *domain.ClassEnrollment) error { return nil }
func (m *MockRepoForTeacher) GetStudentsByClassID(classID uint, year string) ([]domain.Student, error) {
	args := m.Called(classID, year)
	return args.Get(0).([]domain.Student), args.Error(1)
}

// Subject
func (m *MockRepoForTeacher) CreateSubject(subject *domain.Subject) error     { return nil }
func (m *MockRepoForTeacher) GetSubjectByID(id uint) (*domain.Subject, error) { return nil, nil }
func (m *MockRepoForTeacher) AssignSubjectToClass(assignment *domain.ClassSubjectAssignment) error {
	return nil
}
func (m *MockRepoForTeacher) GetAssignmentsByTeacherID(teacherID uint) ([]domain.ClassSubjectAssignment, error) {
	args := m.Called(teacherID)
	return args.Get(0).([]domain.ClassSubjectAssignment), args.Error(1)
}

// Mark
func (m *MockRepoForTeacher) CreateMark(mark *domain.Mark) error {
	args := m.Called(mark)
	return args.Error(0)
}
func (m *MockRepoForTeacher) GetMarksByStudentID(studentID uint, classID uint, subjectID uint) ([]domain.Mark, error) {
	args := m.Called(studentID, classID, subjectID)
	return args.Get(0).([]domain.Mark), args.Error(1)
}
func (m *MockRepoForTeacher) GetMarksByClassID(classID uint) ([]domain.Mark, error) { return nil, nil }
func (m *MockRepoForTeacher) UpdateMark(mark *domain.Mark) error                    { return nil }

// Absence
func (m *MockRepoForTeacher) CreateAbsence(absence *domain.Absence) error { return nil }
func (m *MockRepoForTeacher) GetAbsencesByStudentID(studentID uint, year string) ([]domain.Absence, error) {
	return nil, nil
}
func (m *MockRepoForTeacher) GetAbsencesByClassID(classID uint, date time.Time) ([]domain.Absence, error) {
	return nil, nil
}
func (m *MockRepoForTeacher) UpdateAbsence(absence *domain.Absence) error { return nil }

type MockUserRepoForTeacher struct {
	mock.Mock
}

func (m *MockUserRepoForTeacher) FindByID(id uint) (*domain.User, error) { return nil, nil }

// ... (Add others for domain.UserRepository)
func (m *MockUserRepoForTeacher) Create(u *domain.User) error                    { return nil }
func (m *MockUserRepoForTeacher) FindByEmail(email string) (*domain.User, error) { return nil, nil }
func (m *MockUserRepoForTeacher) Update(u *domain.User) error                    { return nil }
func (m *MockUserRepoForTeacher) FindBySchoolID(schoolID uint) ([]domain.User, error) {
	return nil, nil
}
func (m *MockUserRepoForTeacher) Delete(id uint) error { return nil }
func (m *MockUserRepoForTeacher) GetByExternalID(ctx context.Context, externalID string) (*domain.User, error) {
	return nil, nil
}
func (m *MockUserRepoForTeacher) FindAll(schoolID uint) ([]domain.User, error) { return nil, nil }

func TestCreateMark(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mocks
	mockAcadRepo := new(MockRepoForTeacher)
	mockUserRepo := new(MockUserRepoForTeacher)
	// mockBroadcaster := ... (complex)

	// This is getting complicated to mock entire service dependencies just for a handler test.
	// Instead maybe test the parsing/validation only?
	// Or simpler: Just rely on the compilation test for now if mocks are too heavy?
	// User asked explicit to "Add tests...".

	// Let's implement a simplified test that at least calls the handler method
	// even if it mocks the repository call locally.

	mockAcadRepo.On("CreateMark", mock.AnythingOfType("*domain.Mark")).Return(nil)

	svc := academic.NewAcademicService(mockAcadRepo, mockUserRepo, nil)
	h := NewTeacherHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", uint(1)) // Auth middleware mock

	mark := domain.Mark{
		StudentID: 10,
		SubjectID: 5,
		Value:     8.5,
	}
	body, _ := json.Marshal(mark)
	c.Request, _ = http.NewRequest("POST", "/teacher/marks", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	h.CreateMark(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockAcadRepo.AssertExpectations(t)
}

func TestGetClasses(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAcadRepo := new(MockRepoForTeacher)
	mockUserRepo := new(MockUserRepoForTeacher)

	// Mock data
	assignments := []domain.ClassSubjectAssignment{
		{
			ClassID:   1,
			SubjectID: 1,
			Class: &domain.Class{
				ID:      1,
				Grade:   1,
				Section: "A",
			},
			Subject: &domain.Subject{
				ID:   1,
				Name: "Math",
			},
		},
	}

	mockAcadRepo.On("GetAssignmentsByTeacherID", uint(1)).Return(assignments, nil)

	svc := academic.NewAcademicService(mockAcadRepo, mockUserRepo, nil)
	h := NewTeacherHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", uint(1))

	c.Request, _ = http.NewRequest("GET", "/teacher/classes", nil)

	h.GetClasses(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []domain.ClassSubjectAssignment
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, "A", response[0].Class.Section)

	mockAcadRepo.AssertExpectations(t)
}

func TestGetStudents(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAcadRepo := new(MockRepoForTeacher)
	mockUserRepo := new(MockUserRepoForTeacher)

	students := []domain.Student{
		{ID: 10, FirstName: "Mario", LastName: "Rossi"},
	}

	// Update MockRepoForTeacher to support GetStudentsByClassID with m.Called if not already
	// Since I cannot verify previous edit applied efficiently here without checking, I'll rely on "Update mock methods" part being correct or fix it if it fails compilation.
	// Actually I missed updating GetStudentsByClassID in previous chunk. I should do it now or in separate call.
	// I'll assume I need to update it here or it will panic/return nil.
	// The previous chunk only updated GetAssignments... and GetMarks...
	// I will update the mock method here in a separate block if I could, but I can't do multiple replace calls in one go easily if they overlap or are far apart.
	// I will add the test assuming the mock is updated, then I will update the mock in another call if needed, OR I will try to update it in the previous call if I can edit the tool call? No I can't.
	// I will assume I need to update `GetStudentsByClassID` in the mock.

	mockAcadRepo.On("GetStudentsByClassID", uint(1), "2024-25").Return(students, nil)

	svc := academic.NewAcademicService(mockAcadRepo, mockUserRepo, nil)
	h := NewTeacherHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "classId", Value: "1"}}

	c.Request, _ = http.NewRequest("GET", "/teacher/classes/1/students", nil)

	h.GetStudents(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockAcadRepo.AssertExpectations(t)
}
