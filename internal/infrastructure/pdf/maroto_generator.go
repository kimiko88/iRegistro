package pdf

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/k/iRegistro/internal/domain"
)

type MarotoGenerator struct{}

func NewMarotoGenerator() *MarotoGenerator {
	return &MarotoGenerator{}
}

func (g *MarotoGenerator) GenerateReportCard(data domain.JSONMap) ([]byte, error) {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(20, 10, 20)

	// Header
	m.RegisterHeader(func() {
		m.Row(20, func() {
			m.Col(12, func() {
				m.Text("iRegistro - Official Report Card", props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  16,
					Align: consts.Center,
				})
			})
		})
	})

	// Footer
	m.RegisterFooter(func() {
		m.Row(10, func() {
			m.Col(12, func() {
				m.Text(fmt.Sprintf("Generated on %s", time.Now().Format("2006-01-02 15:04")), props.Text{
					Top:   5,
					Size:  8,
					Align: consts.Right,
				})
			})
		})
	})

	// Content
	// Parse Data
	// Expecting JSON structure for student, stats, marks
	// For now simple dump of JSON map key-values or structured access if possible

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Academic Performance Report", props.Text{
				Size:  14,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
	})

	m.Line(1.0)

	// Dynamically print map content for now (generic template)
	// In real world, we would map specific fields like "Student Name", "Class", "Grades"

	for k, v := range data {
		m.Row(10, func() {
			m.Col(4, func() {
				m.Text(k, props.Text{Style: consts.Bold})
			})
			m.Col(8, func() {
				valStr := ""
				// Simple string conversion
				if str, ok := v.(string); ok {
					valStr = str
				} else {
					b, _ := json.Marshal(v)
					valStr = string(b)
				}
				m.Text(valStr, props.Text{})
			})
		})
	}

	buff, err := m.Output()
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
