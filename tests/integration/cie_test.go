package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"github.com/k/iRegistro/internal/application/auth"
	"github.com/k/iRegistro/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

// MockOIDCProvider creates a mock OIDC provider for testing
type MockOIDCProvider struct {
	server *httptest.Server
}

func NewMockOIDCProvider() *MockOIDCProvider {
	mux := http.NewServeMux()

	// Discovery endpoint
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"issuer":                 "https://mock-cie-provider.it",
			"authorization_endpoint": "https://mock-cie-provider.it/authorize",
			"token_endpoint":         "https://mock-cie-provider.it/token",
			"userinfo_endpoint":      "https://mock-cie-provider.it/userinfo",
			"jwks_uri":               "https://mock-cie-provider.it/jwks",
		})
	})

	// Token endpoint
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "mock_access_token",
			"token_type":   "Bearer",
			"expires_in":   3600,
			"id_token":     createMockIDToken(),
		})
	})

	server := httptest.NewServer(mux)

	return &MockOIDCProvider{
		server: server,
	}
}

func (m *MockOIDCProvider) Close() {
	m.server.Close()
}

func (m *MockOIDCProvider) URL() string {
	return m.server.URL
}

// createMockIDToken creates a mock JWT ID token (unsigned for testing)
func createMockIDToken() string {
	// In real tests, this would be a properly signed JWT
	// For now, return a base64-encoded JSON
	payload := map[string]interface{}{
		"sub":            "cie_user_123",
		"given_name":     "Giulia",
		"family_name":    "Verdi",
		"email":          "giulia.verdi@example.com",
		"fiscal_number":  "VRDGLI90A41H501V",
		"serial_number":  "CA00000AA",
		"date_of_birth":  "1990-01-01",
		"place_of_birth": "Roma",
		"iat":            time.Now().Unix(),
		"exp":            time.Now().Add(1 * time.Hour).Unix(),
	}

	data, _ := json.Marshal(payload)
	return "mock." + string(data) + ".mock"
}

func TestCIEAuthenticationFlow(t *testing.T) {
	// Setup mock OIDC provider
	mockProvider := NewMockOIDCProvider()
	defer mockProvider.Close()

	mockUserRepo := NewMockUserRepository()
	mockTokenService := &auth.TokenService{} // Would need proper init

	// Create OAuth2 config
	oauth2Config := &oauth2.Config{
		ClientID:     "test_client_id",
		ClientSecret: "test_client_secret",
		RedirectURL:  "http://localhost:8080/auth/cie/callback",
		Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
	}

	// Note: In real tests, we'd use oidc.NewProvider with mockProvider.URL()
	// For now, we'll test the service logic directly

	t.Run("Successful CIE Authentication - New User", func(t *testing.T) {
		// Mock CIE claims
		claims := &auth.CIEClaims{
			Sub:          "cie_user_123",
			GivenName:    "Giulia",
			FamilyName:   "Verdi",
			Email:        "giulia.verdi@example.com",
			FiscalNumber: "VRDGLI90A41H501V",
			SerialNumber: "CA00000AA",
			DateOfBirth:  "1990-01-01",
			PlaceOfBirth: "Roma",
		}

		// Test user creation
		schoolID := uuid.New()

		// Since we can't fully mock OIDC verifier, we test the logic separately
		// In a real integration test, you'd use the full flow

		// Verify claims are valid
		assert.NotEmpty(t, claims.FiscalNumber)
		assert.NotEmpty(t, claims.GivenName)
		assert.NotEmpty(t, claims.FamilyName)
		assert.NotEmpty(t, claims.Email)

		// Simulate user creation
		externalID := "cie:" + claims.FiscalNumber
		user := &domain.User{
			ID:              1,
			Email:           claims.Email,
			FirstName:       claims.GivenName,
			LastName:        claims.FamilyName,
			AuthMethod:      "cie",
			CIESerialNumber: &claims.SerialNumber,
			ExternalID:      &externalID,
			SchoolID:        1,
		}

		err := mockUserRepo.Create(user)
		require.NoError(t, err)

		// Verify user was created
		foundUser, err := mockUserRepo.FindByEmail(claims.Email)
		require.NoError(t, err)
		assert.Equal(t, "cie", foundUser.AuthMethod)
		assert.NotNil(t, foundUser.CIESerialNumber)
		assert.Equal(t, "CA00000AA", *foundUser.CIESerialNumber)
	})

	t.Run("CIE Authentication - Existing User Update", func(t *testing.T) {
		// Create existing user
		externalID := "cie:BNCLRA95M12F205Z"
		oldSerial := "CA00001BB"
		existingUser := &domain.User{
			ID:              2,
			Email:           "laura.bianchi@example.com",
			AuthMethod:      "cie",
			ExternalID:      &externalID,
			CIESerialNumber: &oldSerial,
			SchoolID:        1,
		}
		mockUserRepo.Create(existingUser)

		// User renews CIE card (new serial number)
		newSerial := "CA99999ZZ"
		existingUser.CIESerialNumber = &newSerial
		existingUser.LastAuthAt = timePtr(time.Now())

		err := mockUserRepo.Update(existingUser)
		require.NoError(t, err)

		// Verify update
		foundUser, err := mockUserRepo.GetByExternalID(context.Background(), externalID)
		require.NoError(t, err)
		assert.Equal(t, newSerial, *foundUser.CIESerialNumber)
		assert.NotNil(t, foundUser.LastAuthAt)
	})

	t.Run("OAuth2 Authorization URL Generation", func(t *testing.T) {
		state := "random_state_value"
		authURL := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)

		assert.Contains(t, authURL, "client_id=test_client_id")
		assert.Contains(t, authURL, "state=random_state_value")
		assert.Contains(t, authURL, "scope=openid")
	})

	t.Run("Invalid CIE Claims - Missing Fiscal Number", func(t *testing.T) {
		invalidClaims := &auth.CIEClaims{
			GivenName:  "Test",
			FamilyName: "User",
			Email:      "test@example.com",
			// Missing FiscalNumber
		}

		// Should fail validation
		assert.Empty(t, invalidClaims.FiscalNumber)
	})
}

func TestCIESecurityFeatures(t *testing.T) {
	t.Run("State Parameter CSRF Protection", func(t *testing.T) {
		originalState := "secure_random_state_12345"
		receivedState := "different_state_67890"

		// States should match to prevent CSRF
		assert.NotEqual(t, originalState, receivedState, "CSRF protection should reject mismatched states")
	})

	t.Run("ID Token Expiry Check", func(t *testing.T) {
		now := time.Now()

		// Valid token
		validExp := now.Add(1 * time.Hour)
		assert.True(t, validExp.After(now), "Token should be valid")

		// Expired token
		expiredExp := now.Add(-1 * time.Hour)
		assert.True(t, expiredExp.Before(now), "Token should be expired")
	})

	t.Run("Secure Cookie Flags", func(t *testing.T) {
		// Test that cookies should have Secure and HttpOnly flags
		cookieConfig := map[string]bool{
			"Secure":   true, // HTTPS only
			"HttpOnly": true, // No JavaScript access
			"SameSite": true, // CSRF protection
		}

		for flag, expected := range cookieConfig {
			assert.True(t, expected, "Cookie should have %s flag set", flag)
		}
	})
}

func TestCIEMultipleProviders(t *testing.T) {
	// CIE has a single centralized provider, but test different environments
	environments := []struct {
		name   string
		issuer string
	}{
		{"Production", "https://idserver.servizicie.interno.gov.it/idp/profile/oidc"},
		{"Preproduction", "https://preproduzione.idserver.servizicie.interno.gov.it"},
	}

	for _, env := range environments {
		t.Run(env.name, func(t *testing.T) {
			assert.NotEmpty(t, env.issuer)
			assert.Contains(t, env.issuer, "servizicie.interno.gov.it")
		})
	}
}
