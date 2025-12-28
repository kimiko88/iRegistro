package persistence

import (
	"time"

	"github.com/k/iRegistro/internal/domain"
	"gorm.io/gorm"
)

type AcademicRepository struct {
	db *gorm.DB
}

func NewAcademicRepository(db *gorm.DB) *AcademicRepository {
	return &AcademicRepository{db: db}
}

// --- School/Campus ---

func (r *AcademicRepository) CreateSchool(school *domain.School) error {
	return r.db.Create(school).Error
}

func (r *AcademicRepository) GetSchoolByID(id uint) (*domain.School, error) {
	var school domain.School
	if err := r.db.First(&school, id).Error; err != nil {
		return nil, err
	}
	return &school, nil
}

func (r *AcademicRepository) CreateCampus(campus *domain.Campus) error {
	return r.db.Create(campus).Error
}

func (r *AcademicRepository) GetCampusesBySchoolID(schoolID uint) ([]domain.Campus, error) {
	var campuses []domain.Campus
	err := r.db.Where("school_id = ?", schoolID).Find(&campuses).Error
	return campuses, err
}

// --- Curriculum/Class ---

func (r *AcademicRepository) CreateCurriculum(curriculum *domain.Curriculum) error {
	return r.db.Create(curriculum).Error
}

func (r *AcademicRepository) GetCurriculumsBySchoolID(schoolID uint) ([]domain.Curriculum, error) {
	// Join with Campus to filter by SchoolID
	var curriculums []domain.Curriculum
	err := r.db.Joins("JOIN campuses ON campuses.id = curriculums.campus_id").
		Where("campuses.school_id = ?", schoolID).
		Find(&curriculums).Error
	return curriculums, err
}

func (r *AcademicRepository) CreateClass(class *domain.Class) error {
	return r.db.Create(class).Error
}

func (r *AcademicRepository) GetClassByID(id uint) (*domain.Class, error) {
	var class domain.Class
	if err := r.db.First(&class, id).Error; err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *AcademicRepository) GetClassesBySchoolID(schoolID uint) ([]domain.Class, error) {
	var classes []domain.Class
	// Complex join: Class -> Curriculum -> Campus -> School
	err := r.db.Joins("JOIN curriculums ON curriculums.id = classes.curriculum_id").
		Joins("JOIN campuses ON campuses.id = curriculums.campus_id").
		Where("campuses.school_id = ?", schoolID).
		Find(&classes).Error
	return classes, err
}

// --- Student ---

func (r *AcademicRepository) CreateStudent(student *domain.Student) error {
	return r.db.Create(student).Error
}

func (r *AcademicRepository) GetStudentByID(id uint) (*domain.Student, error) {
	var student domain.Student
	if err := r.db.First(&student, id).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *AcademicRepository) EnrollStudent(enrollment *domain.ClassEnrollment) error {
	return r.db.Create(enrollment).Error
}

func (r *AcademicRepository) GetStudentsByClassID(classID uint, year string) ([]domain.Student, error) {
	var students []domain.Student
	err := r.db.Joins("JOIN class_enrollments ON class_enrollments.student_id = students.id").
		Where("class_enrollments.class_id = ? AND class_enrollments.year = ? AND class_enrollments.status = ?", classID, year, domain.EnrollmentActive).
		Find(&students).Error
	return students, err
}

// --- Subject ---

func (r *AcademicRepository) CreateSubject(subject *domain.Subject) error {
	return r.db.Create(subject).Error
}

func (r *AcademicRepository) GetSubjectByID(id uint) (*domain.Subject, error) {
	var subject domain.Subject
	if err := r.db.First(&subject, id).Error; err != nil {
		return nil, err
	}
	return &subject, nil
}

func (r *AcademicRepository) AssignSubjectToClass(assignment *domain.ClassSubjectAssignment) error {
	return r.db.Create(assignment).Error
}

func (r *AcademicRepository) GetAssignmentsByTeacherID(teacherID uint) ([]domain.ClassSubjectAssignment, error) {
	var assignments []domain.ClassSubjectAssignment
	err := r.db.Preload("Class").Preload("Subject").Where("teacher_id = ?", teacherID).Find(&assignments).Error
	return assignments, err
}

// --- Mark ---

func (r *AcademicRepository) CreateMark(mark *domain.Mark) error {
	return r.db.Create(mark).Error
}

func (r *AcademicRepository) GetMarksByStudentID(studentID uint, classID uint, subjectID uint) ([]domain.Mark, error) {
	var marks []domain.Mark
	query := r.db.Where("student_id = ? AND class_id = ?", studentID, classID)
	if subjectID != 0 {
		query = query.Where("subject_id = ?", subjectID)
	}
	err := query.Find(&marks).Error
	return marks, err
}

func (r *AcademicRepository) GetMarksByClassID(classID uint) ([]domain.Mark, error) {
	var marks []domain.Mark
	err := r.db.Where("class_id = ?", classID).Find(&marks).Error
	return marks, err
}

func (r *AcademicRepository) UpdateMark(mark *domain.Mark) error {
	return r.db.Save(mark).Error
}

// --- Absence ---

func (r *AcademicRepository) CreateAbsence(absence *domain.Absence) error {
	return r.db.Create(absence).Error
}

func (r *AcademicRepository) GetAbsencesByStudentID(studentID uint, year string) ([]domain.Absence, error) {
	// TODO: Filter by year (requires joining class enrollment or date range logic)
	// For now, filtering by studentID only or date range logic in service
	var absences []domain.Absence
	err := r.db.Where("student_id = ?", studentID).Find(&absences).Error
	return absences, err
}

func (r *AcademicRepository) GetAbsencesByClassID(classID uint, date time.Time) ([]domain.Absence, error) {
	// Get absences for a specific day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var absences []domain.Absence
	err := r.db.Where("class_id = ? AND date >= ? AND date < ?", classID, startOfDay, endOfDay).Find(&absences).Error
	return absences, err
}

func (r *AcademicRepository) UpdateAbsence(absence *domain.Absence) error {
	return r.db.Save(absence).Error
}
