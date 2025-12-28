package graphql

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

import (
	"github.com/k/iRegistro/internal/application/academic"
)

type Resolver struct {
	AcademicService *academic.AcademicService
}
