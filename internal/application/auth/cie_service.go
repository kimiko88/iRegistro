package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/k/iRegistro/internal/domain"
	"golang.org/x/oauth2"
)

var (
	ErrCIETokenInvalid = errors.New("CIE token validation failed")
	ErrCIEClaimMissing = errors.New("required CIE claim missing")
)

// CIEService handles CIE (Carta d'Identità Elettronica) authentication
type CIEService struct {
	userRepo     domain.UserRepository
	jwtSecret    string
	oauth2Config *oauth2.Config
	oidcVerifier *oidc.IDTokenVerifier
}

// CIEClaims represents claims from CIE OIDC token
type CIEClaims struct {
	Sub          string `json:"sub"`
	Name         string `json:"name"`
	FamilyName   string `json:"family_name"`
	GivenName    string `json:"given_name"`
	Email        string `json:"email"`
	FiscalNumber string `json:"fiscal_number"` // Codice Fiscale
	SerialNumber string `json:"serial_number"` // CIE serial number
	DateOfBirth  string `json:"date_of_birth"`
	PlaceOfBirth string `json:"place_of_birth"`
}

func NewCIEService(userRepo domain.UserRepository, jwtSecret string, oauth2Config *oauth2.Config, verifier *oidc.IDTokenVerifier) *CIEService {
	return &CIEService{
		userRepo:     userRepo,
		jwtSecret:    jwtSecret,
		oauth2Config: oauth2Config,
		oidcVerifier: verifier,
	}
}

// ExchangeCodeForToken exchanges authorization code for OAuth2 token
func (s *CIEService) ExchangeCodeForToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := s.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// ValidateAndExtractClaims validates ID token and extracts CIE claims
func (s *CIEService) ValidateAndExtractClaims(ctx context.Context, rawIDToken string) (*CIEClaims, error) {
	// Verify ID token
	idToken, err := s.oidcVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, ErrCIETokenInvalid
	}

	// Extract claims
	var claims CIEClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}

	// Validate required claims
	if claims.FiscalNumber == "" {
		return nil, ErrCIEClaimMissing
	}
	if claims.GivenName == "" || claims.FamilyName == "" {
		return nil, ErrCIEClaimMissing
	}

	return &claims, nil
}

// GetOrCreateUserByCIE finds or creates a user based on CIE claims
func (s *CIEService) GetOrCreateUserByCIE(ctx context.Context, claims *CIEClaims, schoolID uint) (*domain.User, error) {
	// Try to find user by external ID (CIE unique identifier)
	externalID := "cie:" + claims.FiscalNumber
	user, err := s.userRepo.GetByExternalID(ctx, externalID)

	if err == nil {
		// User exists, update last auth and CIE info
		user.LastAuthAt = timePtr(time.Now())
		user.CIESerialNumber = &claims.SerialNumber
		if err := s.userRepo.Update(user); err != nil {
			return nil, err
		}
		return user, nil
	}

	// User doesn't exist, create new one
	user = &domain.User{
		Email:           claims.Email,
		SchoolID:        schoolID,
		FirstName:       claims.GivenName,
		LastName:        claims.FamilyName,
		AuthMethod:      "cie",
		CIESerialNumber: &claims.SerialNumber,
		ExternalID:      &externalID,
		LastAuthAt:      timePtr(time.Now()),
		Role:            domain.RoleParent, // Default role
		PasswordHash:    "",                // No password for CIE users
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GenerateTokenForUser generates JWT token for authenticated CIE user
func (s *CIEService) GenerateTokenForUser(user *domain.User) (string, error) {
	// Simple JWT generation - in production use a proper token service
	return fmt.Sprintf("jwt_token_user_%d_school_%d", user.ID, user.SchoolID), nil
}

// GetAuthorizationURL generates OAuth2 authorization URL for CIE login
func (s *CIEService) GetAuthorizationURL(state string) string {
	return s.oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}
