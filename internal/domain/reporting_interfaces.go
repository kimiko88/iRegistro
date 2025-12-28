package domain

import "time"

type ReportingRepository interface {
	// Documents
	CreateDocument(doc *Document) error
	GetDocumentByID(id uint) (*Document, error)
	GetDocumentsBySchoolID(schoolID uint, docType DocumentType) ([]Document, error)
	GetDocumentsByStatus(schoolID uint, status DocumentStatus) ([]Document, error)
	GetDocumentsByStudentID(studentID uint) ([]Document, error)
	UpdateDocument(doc *Document) error
	DeleteDocument(id uint) error

	// Signatures
	AddSignature(sig *DocumentSignature) error
	GetSignaturesByDocumentID(docID uint) ([]DocumentSignature, error)

	// PCTO
	CreatePCTOProject(project *PCTOProject) error
	GetPCTOProjectsBySchoolID(schoolID uint) ([]PCTOProject, error)
	CreatePCTOAssignment(assignment *PCTOAssignment) error
	GetPCTOAssignmentsByClassID(classID uint) ([]PCTOAssignment, error)
	GetPCTOAssignmentsByStudentID(studentID uint) ([]PCTOAssignment, error)
	LogPCTOHours(hours *PCTOHour) error

	// Orientation
	CreateOrientationActivity(activity *OrientationActivity) error
	GetOrientationActivitiesBySchoolID(schoolID uint) ([]OrientationActivity, error)
	RegisterOrientationParticipation(participation *OrientationParticipation) error
	GetOrientationParticipationsByStudentID(studentID uint) ([]OrientationParticipation, error)

	// Stats
	CountDocumentsByStatus(schoolID uint, status DocumentStatus) (int64, error)
	CountDocumentsUpdatedSince(schoolID uint, status []DocumentStatus, since time.Time) (int64, error)
}

type PDFGenerator interface {
	GenerateReportCard(data JSONMap) ([]byte, error)
}
