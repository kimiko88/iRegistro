package pdf_test

import (
	"testing"

	"github.com/k/iRegistro/internal/domain"
	"github.com/k/iRegistro/internal/infrastructure/pdf"
	"github.com/stretchr/testify/assert"
)

func TestGenerateReportCard(t *testing.T) {
	gen := pdf.NewMarotoGenerator()

	data := domain.JSONMap{
		"Student": "Mario Rossi",
		"Class":   "1A",
		"Math":    "8",
		"History": "7",
	}

	pdfBytes, err := gen.GenerateReportCard(data)
	assert.NoError(t, err)
	assert.NotEmpty(t, pdfBytes)
	assert.Contains(t, string(pdfBytes), "%PDF")
}
