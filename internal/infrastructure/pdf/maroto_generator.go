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

	g.addHeader(m, "Official Report Card")

	// Content
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Academic Performance Report", props.Text{
				Size:  14,
				Style: consts.Bold,
				Align: consts.Left,
			})
		})
	})

	g.addKeyValueTable(m, data)
	g.addFooter(m)

	buff, err := m.Output()
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func (g *MarotoGenerator) GenerateCertificate(data domain.JSONMap) ([]byte, error) {
	m := pdf.NewMaroto(consts.Landscape, consts.A4)
	m.SetPageMargins(30, 20, 30)

	// Header
	m.RegisterHeader(func() {
		m.Row(40, func() {
			m.Col(12, func() {
				m.Text("CERTIFICATE OF ATTENDANCE", props.Text{
					Top:   15,
					Style: consts.BoldItalic,
					Size:  24,
					Align: consts.Center,
				})
			})
		})
	})

	// Content
	m.Row(20, func() {
		m.Col(12, func() {
			m.Text("This verifies that", props.Text{Align: consts.Center, Size: 12})
		})
	})

	studentName, _ := data["student_name"].(string)
	m.Row(15, func() {
		m.Col(12, func() {
			m.Text(studentName, props.Text{Align: consts.Center, Size: 18, Style: consts.Bold})
		})
	})

	g.addFooter(m)

	buff, err := m.Output()
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// --- Helpers ---

func (g *MarotoGenerator) addHeader(m pdf.Maroto, title string) {
	m.RegisterHeader(func() {
		m.Row(20, func() {
			m.Col(12, func() {
				m.Text("iRegistro - "+title, props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  16,
					Align: consts.Center,
				})
			})
		})
	})
}

func (g *MarotoGenerator) addFooter(m pdf.Maroto) {
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
}

func (g *MarotoGenerator) addKeyValueTable(m pdf.Maroto, data domain.JSONMap) {
	m.Line(1.0)
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
}
