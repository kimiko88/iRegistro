package domain

import "time"

type AcademicRepository interface {
	// School/Campus
	CreateSchool(school *School) error
	GetSchoolByID(id uint) (*School, error)
	GetAllSchools() ([]School, error)
	CountSchools() (int64, error)
	CreateCampus(campus *Campus) error
	GetCampusesBySchoolID(schoolID uint) ([]Campus, error)

	// Curriculum/Class
	CreateCurriculum(curriculum *Curriculum) error
	GetCurriculumsBySchoolID(schoolID uint) ([]Curriculum, error)
	CreateClass(class *Class) error
	GetClassByID(id uint) (*Class, error)
	GetClassesBySchoolID(schoolID uint) ([]Class, error)

	// Student
	CreateStudent(student *Student) error
	GetStudentByID(id uint) (*Student, error)
	EnrollStudent(enrollment *ClassEnrollment) error
	GetStudentsByClassID(classID uint, year string) ([]Student, error)

	// Subject
	CreateSubject(subject *Subject) error
	GetSubjectByID(id uint) (*Subject, error)
	GetSubjectsByIDs(ids []uint) ([]Subject, error) // Added
	AssignSubjectToClass(assignment *ClassSubjectAssignment) error
	GetAssignmentsByTeacherID(teacherID uint) ([]ClassSubjectAssignment, error) // Added

	// Mark
	CreateMark(mark *Mark) error
	GetMarksByStudentID(studentID uint, classID uint, subjectID uint) ([]Mark, error)
	GetMarksByClassID(classID uint) ([]Mark, error) // For averages
	UpdateMark(mark *Mark) error

	// Absence
	CreateAbsence(absence *Absence) error
	GetAbsencesByStudentID(studentID uint, year string) ([]Absence, error)
	GetAbsencesByClassID(classID uint, date time.Time) ([]Absence, error)
	UpdateAbsence(absence *Absence) error
}

type AcademicService interface {
	// ... (To be implemented)
	CalculateAverage(marks []Mark) float64
	CalculateWeightedAverage(marks []Mark) float64
	CheckAbsenceThreshold(studentID uint, totalHours int) (bool, float64, error)
}
