package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/k/iRegistro/internal/application/academic"
	"github.com/k/iRegistro/internal/domain"
)

type AcademicHandler struct {
	service *academic.AcademicService
}

func NewAcademicHandler(service *academic.AcademicService) *AcademicHandler {
	return &AcademicHandler{service: service}
}

// --- Schools ---

func (h *AcademicHandler) CreateSchool(c *gin.Context) {
	var school domain.School
	if err := c.ShouldBindJSON(&school); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateSchool(&school); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, school)
}

// --- Campuses ---

func (h *AcademicHandler) GetCampuses(c *gin.Context) {
	schoolID, err := strconv.Atoi(c.Param("schoolId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid school id"})
		return
	}
	campuses, err := h.service.GetCampusesBySchoolID(uint(schoolID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, campuses)
}

func (h *AcademicHandler) CreateCampus(c *gin.Context) {
	schoolID, err := strconv.Atoi(c.Param("schoolId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid school id"})
		return
	}
	var campus domain.Campus
	if err := c.ShouldBindJSON(&campus); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	campus.SchoolID = uint(schoolID)
	if err := h.service.CreateCampus(&campus); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, campus)
}

// --- Curriculums ---

func (h *AcademicHandler) GetCurriculums(c *gin.Context) {
	schoolID, err := strconv.Atoi(c.Param("schoolId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid school id"})
		return
	}
	curriculums, err := h.service.GetCurriculumsBySchoolID(uint(schoolID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, curriculums)
}

func (h *AcademicHandler) CreateCurriculum(c *gin.Context) {
	// ... similar pattern
	var curr domain.Curriculum
	if err := c.ShouldBindJSON(&curr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateCurriculum(&curr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, curr)
}

// --- Subjects ---

func (h *AcademicHandler) GetSubjects(c *gin.Context) {
	schoolIDStr := c.Param("schoolId")
	schoolID, err := strconv.ParseUint(schoolIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid school ID"})
		return
	}

	subjects, err := h.service.GetSubjectsBySchool(uint(schoolID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subjects)
}

func (h *AcademicHandler) CreateSubject(c *gin.Context) {
	schoolIDStr := c.Param("schoolId")
	schoolID, err := strconv.ParseUint(schoolIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid school ID"})
		return
	}

	var subject domain.Subject
	if err := c.ShouldBindJSON(&subject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subject.SchoolID = uint(schoolID)

	if err := h.service.CreateSubject(&subject); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, subject)
}

// --- Classes ---

func (h *AcademicHandler) GetClasses(c *gin.Context) {
	schoolID, err := strconv.Atoi(c.Param("schoolId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid school id"})
		return
	}
	classes, err := h.service.GetClassesBySchoolID(uint(schoolID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, classes)
}

func (h *AcademicHandler) CreateClass(c *gin.Context) {
	schoolID, err := strconv.Atoi(c.Param("schoolId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid school id"})
		return
	}

	var class domain.Class
	if err := c.ShouldBindJSON(&class); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	class.SchoolID = uint(schoolID)

	if err := h.service.CreateClass(&class); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, class)
}

func (h *AcademicHandler) GetClassDetails(c *gin.Context) {
	classID, err := strconv.Atoi(c.Param("classId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid class id"})
		return
	}
	class, err := h.service.GetClassByID(uint(classID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, class)
}

// --- Marks ---

func (h *AcademicHandler) CreateMark(c *gin.Context) {
	var mark domain.Mark
	if err := c.ShouldBindJSON(&mark); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO: Get TeacherID from context
	if err := h.service.CreateMark(&mark); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, mark)
}

// --- Assignments ---

func (h *AcademicHandler) AssignSubjectToClass(c *gin.Context) {
	// schoolID, _ := strconv.Atoi(c.Param("schoolId")) // Not explicitly used but good context
	var assignment domain.ClassSubjectAssignment
	if err := c.ShouldBindJSON(&assignment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Validate SchoolID matches Class/Subject/Teacher

	if err := h.service.AssignSubjectToClass(&assignment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, assignment)
}
