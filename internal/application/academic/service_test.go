package academic

import (
	"testing"

	"github.com/k/iRegistro/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCalculateAverage(t *testing.T) {
	service := NewAcademicService(nil, nil, nil) // Repo not needed for pure logic

	marks := []domain.Mark{
		{Value: 8.0},
		{Value: 6.0},
		{Value: 7.0},
	}
	avg := service.CalculateAverage(marks)
	assert.Equal(t, 7.0, avg)

	// Empty
	assert.Equal(t, 0.0, service.CalculateAverage([]domain.Mark{}))
}

func TestCalculateWeightedAverage(t *testing.T) {
	service := NewAcademicService(nil, nil, nil)

	marks := []domain.Mark{
		{Value: 8.0, Weight: 1.0}, // e.g. Oral
		{Value: 6.0, Weight: 0.5}, // e.g. Homework
	}
	// (8*1 + 6*0.5) / (1 + 0.5) = (8 + 3) / 1.5 = 11 / 1.5 = 7.333...
	avg := service.CalculateWeightedAverage(marks)
	assert.InDelta(t, 7.33, avg, 0.01)

	assert.Equal(t, 0.0, service.CalculateWeightedAverage([]domain.Mark{{Value: 10, Weight: 0}}))
}

func TestCheckAbsenceThreshold(t *testing.T) {
	service := NewAcademicService(nil, nil, nil)

	// Total 100 hours. 31 absences -> 69% attendance -> Fail
	absences := make([]domain.Absence, 31)
	for i := 0; i < 31; i++ {
		absences[i] = domain.Absence{Type: domain.AbsenceFull}
	}

	fail, rate, _ := service.CheckAbsenceThreshold(absences, 100)
	assert.True(t, fail)
	assert.InDelta(t, 0.31, rate, 0.01)

	// Total 100 hours. 20 absences -> 80% attendance -> Pass
	absencesOk := make([]domain.Absence, 20)
	for i := 0; i < 20; i++ {
		absencesOk[i] = domain.Absence{Type: domain.AbsenceFull}
	}
	failOk, _, _ := service.CheckAbsenceThreshold(absencesOk, 100)
	assert.False(t, failOk)
}
