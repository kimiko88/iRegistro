package benchmarks

import (
	"testing"

	"github.com/k/iRegistro/internal/domain"
	// "github.com/k/iRegistro/internal/application/academic"
)

// BenchmarkCalculateAverage tests logic performance
// Run with: go test -bench=. ./tests/benchmarks
func BenchmarkCalculateAverage(b *testing.B) {
	// Setup data
	marks := make([]domain.Mark, 1000)
	for i := 0; i < 1000; i++ {
		marks[i] = domain.Mark{Value: 7.5, Weight: 1.0}
	}

	// service := academic.NewAcademicService(nil, nil, nil) // No repo needed for calc

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// logic-only benchmark
		var sum float64
		for _, m := range marks {
			sum += m.Value
		}
		_ = sum / float64(len(marks))
	}
}
