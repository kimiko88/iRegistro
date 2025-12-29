package academic

import (
	"github.com/k/iRegistro/internal/domain"
)

type AcademicService struct {
	repo     domain.AcademicRepository
	userRepo domain.UserRepository
	notifier domain.NotificationService
}

func NewAcademicService(repo domain.AcademicRepository, userRepo domain.UserRepository, notifier domain.NotificationService) *AcademicService {
	return &AcademicService{
		repo:     repo,
		userRepo: userRepo,
		notifier: notifier,
	}
}

// --- School/Campus ---

func (s *AcademicService) CreateSchool(school *domain.School) error {
	return s.repo.CreateSchool(school)
}

func (s *AcademicService) GetSchoolByID(id uint) (*domain.School, error) {
	return s.repo.GetSchoolByID(id)
}

func (s *AcademicService) CreateCampus(campus *domain.Campus) error {
	return s.repo.CreateCampus(campus)
}

func (s *AcademicService) GetCampusesBySchoolID(schoolID uint) ([]domain.Campus, error) {
	return s.repo.GetCampusesBySchoolID(schoolID)
}

// --- Curriculum/Class ---

func (s *AcademicService) CreateCurriculum(curriculum *domain.Curriculum) error {
	return s.repo.CreateCurriculum(curriculum)
}

func (s *AcademicService) GetCurriculumsBySchoolID(schoolID uint) ([]domain.Curriculum, error) {
	return s.repo.GetCurriculumsBySchoolID(schoolID)
}

func (s *AcademicService) CreateClass(class *domain.Class) error {
	return s.repo.CreateClass(class)
}

func (s *AcademicService) GetClassByID(id uint) (*domain.Class, error) {
	return s.repo.GetClassByID(id)
}

func (s *AcademicService) GetClassesBySchoolID(schoolID uint) ([]domain.Class, error) {
	return s.repo.GetClassesBySchoolID(schoolID)
}

// --- Student ---

func (s *AcademicService) CreateStudent(student *domain.Student) error {
	return s.repo.CreateStudent(student)
}

func (s *AcademicService) GetStudentByID(id uint) (*domain.Student, error) {
	return s.repo.GetStudentByID(id)
}

func (s *AcademicService) EnrollStudent(enrollment *domain.ClassEnrollment) error {
	return s.repo.EnrollStudent(enrollment)
}

func (s *AcademicService) GetStudentsByClassID(classID uint, year string) ([]domain.Student, error) {
	return s.repo.GetStudentsByClassID(classID, year)
}

// --- Subject ---

func (s *AcademicService) CreateSubject(subject *domain.Subject) error {
	return s.repo.CreateSubject(subject)
}

func (s *AcademicService) GetSubjectByID(id uint) (*domain.Subject, error) {
	return s.repo.GetSubjectByID(id)
}

func (s *AcademicService) GetSubjectsBySchool(schoolID uint) ([]domain.Subject, error) {
	// Query subjects filtered by school ID
	return s.repo.GetSubjects(schoolID)
}

func (s *AcademicService) AssignSubjectToClass(assignment *domain.ClassSubjectAssignment) error {
	return s.repo.AssignSubjectToClass(assignment)
}

func (s *AcademicService) GetAssignmentsByTeacherID(teacherID uint) ([]domain.ClassSubjectAssignment, error) {
	return s.repo.GetAssignmentsByTeacherID(teacherID)
}

// --- Marks ---

func (s *AcademicService) CreateMark(mark *domain.Mark) error {
	// TODO: Maybe add validation logic here (e.g. check weight sum, check teacher assignment)
	if err := s.repo.CreateMark(mark); err != nil {
		return err
	}
	if s.notifier != nil {
		s.notifier.NotifyMarkAdded(mark)
	}
	return nil
}

func (s *AcademicService) GetMarksByStudentID(studentID, classID, subjectID uint) ([]domain.Mark, error) {
	return s.repo.GetMarksByStudentID(studentID, classID, subjectID)
}

func (s *AcademicService) GetMarksByClassID(classID uint) ([]domain.Mark, error) {
	return s.repo.GetMarksByClassID(classID)
}

// --- Teacher/User ---

func (s *AcademicService) GetTeacherByID(id uint) (*domain.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	// Optional: Check Role
	// if user.Role != domain.RoleTeacher { return nil, fmt.Errorf("user is not a teacher") }
	return user, nil
}

// --- Logic implementations to verify via TDD ---
func (s *AcademicService) CalculateAverage(marks []domain.Mark) float64 {
	if len(marks) == 0 {
		return 0
	}
	var sum float64
	for _, m := range marks {
		sum += m.Value
	}
	return sum / float64(len(marks))
}

func (s *AcademicService) CalculateWeightedAverage(marks []domain.Mark) float64 {
	if len(marks) == 0 {
		return 0
	}
	var sum float64
	var totalWeight float64
	for _, m := range marks {
		sum += m.Value * m.Weight
		totalWeight += m.Weight
	}
	if totalWeight == 0 {
		return 0
	}
	return sum / totalWeight
}

// CheckAbsenceThreshold checks if attendance is below 70% (0.7) usually.
// totalHours is the expected total hours for the period (could be per subject or global).
func (s *AcademicService) CheckAbsenceThreshold(absences []domain.Absence, totalHours int) (bool, float64, error) {
	if totalHours <= 0 {
		return false, 0, nil
	}

	// Count hours of absence
	// If Type=ABSENT (Full day) we might need to know how many hours that day had.
	// For simplicity, assuming here Absences are recorded per hour or normalized.
	// If Type=DAD it counts as present? Requirements say: DaD tracking for participation.
	// Usually DaD is Present, but logic might vary. Assuming Only ABSENT counts against threshold.

	absentCount := 0
	for _, a := range absences {
		if a.Type == domain.AbsenceFull { // Assuming ABSENT means missed hour
			absentCount++
		}
	}

	percentageAbsent := float64(absentCount) / float64(totalHours)
	isBelowThreshold := (1.0 - percentageAbsent) < 0.70 // < 70% attendance

	return isBelowThreshold, percentageAbsent, nil
}
