package persistence

import (
	"github.com/k/iRegistro/internal/domain"
	"gorm.io/gorm"
)

type ReportingRepository struct {
	db *gorm.DB
}

func NewReportingRepository(db *gorm.DB) *ReportingRepository {
	return &ReportingRepository{db: db}
}

// --- Documents ---

func (r *ReportingRepository) CreateDocument(doc *domain.Document) error {
	return r.db.Create(doc).Error
}

func (r *ReportingRepository) GetDocumentByID(id uint) (*domain.Document, error) {
	var doc domain.Document
	if err := r.db.Preload("Signatures").First(&doc, id).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *ReportingRepository) GetDocumentsBySchoolID(schoolID uint, docType domain.DocumentType) ([]domain.Document, error) {
	var docs []domain.Document
	query := r.db.Where("school_id = ?", schoolID)
	if docType != "" {
		query = query.Where("type = ?", docType)
	}
	err := query.Find(&docs).Error
	return docs, err
}

func (r *ReportingRepository) GetDocumentsByStudentID(studentID uint) ([]domain.Document, error) {
	var docs []domain.Document
	err := r.db.Where("student_id = ?", studentID).Find(&docs).Error
	return docs, err
}

func (r *ReportingRepository) UpdateDocument(doc *domain.Document) error {
	return r.db.Save(doc).Error
}

func (r *ReportingRepository) DeleteDocument(id uint) error {
	return r.db.Delete(&domain.Document{}, id).Error
}

// --- Signatures ---

func (r *ReportingRepository) AddSignature(sig *domain.DocumentSignature) error {
	return r.db.Create(sig).Error
}

func (r *ReportingRepository) GetSignaturesByDocumentID(docID uint) ([]domain.DocumentSignature, error) {
	var sigs []domain.DocumentSignature
	err := r.db.Where("document_id = ?", docID).Find(&sigs).Error
	return sigs, err
}

// --- PCTO ---

func (r *ReportingRepository) CreatePCTOProject(project *domain.PCTOProject) error {
	return r.db.Create(project).Error
}

func (r *ReportingRepository) GetPCTOProjectsBySchoolID(schoolID uint) ([]domain.PCTOProject, error) {
	var projects []domain.PCTOProject
	err := r.db.Where("school_id = ?", schoolID).Find(&projects).Error
	return projects, err
}

func (r *ReportingRepository) CreatePCTOAssignment(assignment *domain.PCTOAssignment) error {
	return r.db.Create(assignment).Error
}

func (r *ReportingRepository) GetPCTOAssignmentsByClassID(classID uint) ([]domain.PCTOAssignment, error) {
	var assignments []domain.PCTOAssignment
	// Join with ClassEnrollment to filter assignments for students in that class?
	// Or simplistic assumption we just iterate students.
	// But assignments table doesn't have ClassID.
	// For efficiency, maybe simpler via Service orchestrating StudentIDs.
	// OR query: join pcto_assignments -> students -> class_enrollments
	err := r.db.Joins("JOIN students ON students.id = pcto_assignments.student_id").
		Joins("JOIN class_enrollments ON class_enrollments.student_id = students.id").
		Where("class_enrollments.class_id = ? AND class_enrollments.status = ?", classID, domain.EnrollmentActive).
		Find(&assignments).Error
	return assignments, err
}

func (r *ReportingRepository) GetPCTOAssignmentsByStudentID(studentID uint) ([]domain.PCTOAssignment, error) {
	var assignments []domain.PCTOAssignment
	err := r.db.Preload("Hours").Where("student_id = ?", studentID).Find(&assignments).Error
	return assignments, err
}

func (r *ReportingRepository) LogPCTOHours(hours *domain.PCTOHour) error {
	return r.db.Create(hours).Error
}

// --- Orientation ---

func (r *ReportingRepository) CreateOrientationActivity(activity *domain.OrientationActivity) error {
	return r.db.Create(activity).Error
}

func (r *ReportingRepository) GetOrientationActivitiesBySchoolID(schoolID uint) ([]domain.OrientationActivity, error) {
	var activities []domain.OrientationActivity
	err := r.db.Where("school_id = ?", schoolID).Find(&activities).Error
	return activities, err
}

func (r *ReportingRepository) RegisterOrientationParticipation(participation *domain.OrientationParticipation) error {
	return r.db.Create(participation).Error
}

func (r *ReportingRepository) GetOrientationParticipationsByStudentID(studentID uint) ([]domain.OrientationParticipation, error) {
	var parts []domain.OrientationParticipation
	err := r.db.Where("student_id = ?", studentID).Find(&parts).Error
	return parts, err
}
