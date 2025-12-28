package reporting

import (
	"errors"
	"time"

	"github.com/k/iRegistro/internal/domain"
)

type ReportingService struct {
	repo   domain.ReportingRepository
	pdfGen domain.PDFGenerator
}

func NewReportingService(repo domain.ReportingRepository, pdfGen domain.PDFGenerator) *ReportingService {
	return &ReportingService{repo: repo, pdfGen: pdfGen}
}

// --- Documents ---

func (s *ReportingService) CreateReportCard(schoolID, studentID, classID uint, academicYear string, data domain.JSONMap, creatorID uint) (*domain.Document, error) {
	doc := &domain.Document{
		SchoolID:     schoolID,
		StudentID:    &studentID,
		ClassID:      &classID,
		Type:         domain.DocReportCard,
		Title:        "Pagella - " + academicYear,
		Data:         data,
		AcademicYear: academicYear,
		CreatedBy:    creatorID,
		CreatedAt:    time.Now(),
		Status:       domain.DocStatusDraft,
	}
	if err := s.repo.CreateDocument(doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *ReportingService) SignDocument(docID, signerID uint, ipAddress string) error {
	doc, err := s.repo.GetDocumentByID(docID)
	if err != nil {
		return err
	}
	if doc.Status == domain.DocStatusArchived {
		return errors.New("cannot sign archived document")
	}

	sig := &domain.DocumentSignature{
		DocumentID:         docID,
		SignerID:           signerID,
		SignatureTimestamp: time.Now(),
		IPAddress:          ipAddress,
		IsValid:            true,
	}

	if err := s.repo.AddSignature(sig); err != nil {
		return err
	}

	// Update status if needed (e.g. if Principal signs, finalized?)
	// For now simple logic: if signed, maybe mark as Signed?
	if doc.Status == domain.DocStatusDraft {
		doc.Status = domain.DocStatusSigned
		return s.repo.UpdateDocument(doc)
	}
	return nil
}

func (s *ReportingService) GetDocumentPDF(docID uint) ([]byte, error) {
	doc, err := s.repo.GetDocumentByID(docID)
	if err != nil {
		return nil, err
	}

	if s.pdfGen == nil {
		return nil, errors.New("pdf generator not configured")
	}

	return s.pdfGen.GenerateReportCard(doc.Data)
}

// ... Implement other service methods wrapping repo ...
func (s *ReportingService) GetDocumentsBySchoolID(schoolID uint, docType domain.DocumentType) ([]domain.Document, error) {
	return s.repo.GetDocumentsBySchoolID(schoolID, docType)
}

// PCTO
func (s *ReportingService) CreatePCTOProject(project *domain.PCTOProject) error {
	project.CreatedAt = time.Now()
	return s.repo.CreatePCTOProject(project)
}
func (s *ReportingService) GetPCTOProjects(schoolID uint) ([]domain.PCTOProject, error) {
	return s.repo.GetPCTOProjectsBySchoolID(schoolID)
}

// Orientation
func (s *ReportingService) RegisterOrientation(participation *domain.OrientationParticipation) error {
	return s.repo.RegisterOrientationParticipation(participation)
}

// Student Reporting
func (s *ReportingService) GetDocumentsByStudentID(studentID uint) ([]domain.Document, error) {
	return s.repo.GetDocumentsByStudentID(studentID)
}

func (s *ReportingService) GetPCTOProgression(studentID uint) (int, []domain.PCTOProject, error) {
	assignments, err := s.repo.GetPCTOAssignmentsByStudentID(studentID)
	if err != nil {
		return 0, nil, err
	}

	totalHours := 0
	var projects []domain.PCTOProject

	// Assuming HoursPlanned is the metric for now
	for _, a := range assignments {
		totalHours += a.HoursPlanned
		// TODO: Fetch project details if needed
	}

	return totalHours, projects, nil
}

func (s *ReportingService) GetOrientationHours(studentID uint) (int, error) {
	participations, err := s.repo.GetOrientationParticipationsByStudentID(studentID)
	if err != nil {
		return 0, err
	}
	total := 0
	for _, p := range participations {
		total += p.HoursEarned
	}
	return total, nil
}
