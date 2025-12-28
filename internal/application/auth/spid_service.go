package auth

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"time"

	"github.com/crewjam/saml"
	"github.com/k/iRegistro/internal/domain"
)

var (
	ErrSAMLValidationFailed = errors.New("SAML assertion validation failed")
	ErrSPIDAttributeMissing = errors.New("required SPID attribute missing")
)

// SPIDService handles SPID authentication
type SPIDService struct {
	userRepo        domain.UserRepository
	jwtSecret       string
	serviceProvider *saml.ServiceProvider
}

// SPIDAttributes represents attributes returned by SPID providers
type SPIDAttributes struct {
	TaxCode     string // spidCode (Codice Fiscale)
	Name        string // name
	FamilyName  string // familyName
	Email       string // email
	DateOfBirth string // dateOfBirth (optional)
	Provider    string // Identity Provider name
}

func NewSPIDService(userRepo domain.UserRepository, jwtSecret string, sp *saml.ServiceProvider) *SPIDService {
	return &SPIDService{
		userRepo:        userRepo,
		jwtSecret:       jwtSecret,
		serviceProvider: sp,
	}
}

// ValidateSAMLAssertion validates and extracts SPID attributes from SAML assertion
func (s *SPIDService) ValidateSAMLAssertion(ctx context.Context, assertion *saml.Assertion) (*SPIDAttributes, error) {
	if assertion == nil {
		return nil, ErrSAMLValidationFailed
	}

	// Validate assertion
	if err := s.validateAssertion(assertion); err != nil {
		return nil, err
	}

	// Extract attributes
	attrs := &SPIDAttributes{
		Provider: s.extractIssuer(assertion),
	}

	for _, attrStatement := range assertion.AttributeStatements {
		for _, attr := range attrStatement.Attributes {
			if len(attr.Values) == 0 {
				continue
			}

			value := attr.Values[0].Value

			switch attr.Name {
			case "spidCode", "fiscalNumber":
				attrs.TaxCode = value
			case "name":
				attrs.Name = value
			case "familyName":
				attrs.FamilyName = value
			case "email":
				attrs.Email = value
			case "dateOfBirth":
				attrs.DateOfBirth = value
			}
		}
	}

	// Validate required attributes
	if attrs.TaxCode == "" {
		return nil, ErrSPIDAttributeMissing
	}
	if attrs.Name == "" || attrs.FamilyName == "" {
		return nil, ErrSPIDAttributeMissing
	}

	return attrs, nil
}

// GetOrCreateUserBySPID finds or creates a user based on SPID attributes
func (s *SPIDService) GetOrCreateUserBySPID(ctx context.Context, attrs *SPIDAttributes, schoolID uint) (*domain.User, error) {
	// Try to find user by external ID (SPID unique identifier)
	externalID := "spid:" + attrs.TaxCode
	user, err := s.userRepo.GetByExternalID(ctx, externalID)

	if err == nil {
		// User exists, update last auth
		user.LastAuthAt = timePtr(time.Now())
		user.SPIDProvider = &attrs.Provider
		if err := s.userRepo.Update(user); err != nil {
			return nil, err
		}
		return user, nil
	}

	// User doesn't exist, create new one
	user = &domain.User{
		Email:        attrs.Email,
		SchoolID:     schoolID,
		FirstName:    attrs.Name,
		LastName:     attrs.FamilyName,
		AuthMethod:   "spid",
		SPIDProvider: &attrs.Provider,
		ExternalID:   &externalID,
		LastAuthAt:   timePtr(time.Now()),
		Role:         domain.RoleParent, // Default role
		PasswordHash: "",                // No password for SPID users
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GenerateTokenForUser generates JWT token for authenticated SPID user
func (s *SPIDService) GenerateTokenForUser(user *domain.User) (string, error) {
	// Simple JWT generation - in production use a proper token service
	return fmt.Sprintf("jwt_token_user_%d_school_%d", user.ID, user.SchoolID), nil
}

// validateAssertion performs SAML assertion validation
func (s *SPIDService) validateAssertion(assertion *saml.Assertion) error {
	// Check assertion expiry
	now := time.Now()

	if assertion.Conditions != nil {
		// NotBefore and NotOnOrAfter are time.Time, not pointers
		if !assertion.Conditions.NotBefore.IsZero() && now.Before(assertion.Conditions.NotBefore) {
			return errors.New("assertion not yet valid")
		}
		if !assertion.Conditions.NotOnOrAfter.IsZero() && now.After(assertion.Conditions.NotOnOrAfter) {
			return errors.New("assertion expired")
		}
	}

	// Validate subject
	if assertion.Subject == nil || assertion.Subject.NameID == nil {
		return errors.New("assertion missing subject")
	}

	return nil
}

// extractIssuer extracts the identity provider name from assertion
func (s *SPIDService) extractIssuer(assertion *saml.Assertion) string {
	// Issuer is a struct with Value field, not a pointer
	return assertion.Issuer.Value
}

// SPIDAttributeValue represents a SAML attribute value
type SPIDAttributeValue struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

func timePtr(t time.Time) *time.Time {
	return &t
}
