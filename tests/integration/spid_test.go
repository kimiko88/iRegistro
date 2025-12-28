package integration

import (
	"context"
	"testing"
	"time"

	"github.com/crewjam/saml"
	"github.com/google/uuid"
	"github.com/k/iRegistro/internal/application/auth"
	"github.com/k/iRegistro/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockUserRepository for testing
type MockUserRepository struct {
	users map[string]*domain.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (m *MockUserRepository) Create(user *domain.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	if user, ok := m.users[email]; ok {
		return user, nil
	}
	return nil, domain.ErrUserNotFound
}

func (m *MockUserRepository) FindByID(id uint) (*domain.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (m *MockUserRepository) GetByExternalID(ctx context.Context, externalID string) (*domain.User, error) {
	for _, user := range m.users {
		if user.ExternalID != nil && *user.ExternalID == externalID {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (m *MockUserRepository) Update(user *domain.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) FindAll(schoolID uint) ([]domain.User, error) {
	return nil, nil
}

// MockAuthRepository for testing
type MockAuthRepository struct {
	sessions      map[string]*domain.Session
	refreshTokens map[string]*domain.RefreshToken
}

func NewMockAuthRepository() *MockAuthRepository {
	return &MockAuthRepository{
		sessions:      make(map[string]*domain.Session),
		refreshTokens: make(map[string]*domain.RefreshToken),
	}
}

func (m *MockAuthRepository) CreateSession(session *domain.Session) error {
	m.sessions[session.TokenHash] = session
	return nil
}

func (m *MockAuthRepository) StoreRefreshToken(token *domain.RefreshToken) error {
	m.refreshTokens[token.TokenHash] = token
	return nil
}

func (m *MockAuthRepository) RevokeRefreshToken(tokenHash string) error {
	if token, ok := m.refreshTokens[tokenHash]; ok {
		now := time.Now()
		token.RevokedAt = &now
	}
	return nil
}

func (m *MockAuthRepository) GetRefreshToken(tokenHash string) (*domain.RefreshToken, error) {
	if token, ok := m.refreshTokens[tokenHash]; ok {
		return token, nil
	}
	return nil, domain.ErrUserNotFound
}

// createMockSAMLAssertion creates a mock SAML assertion for testing
func createMockSAMLAssertion(taxCode, name, familyName, email, provider string) *saml.Assertion {
	now := time.Now()

	assertion := &saml.Assertion{
		ID:           "_" + uuid.New().String(),
		IssueInstant: now,
		Version:      "2.0",
		Issuer: saml.Issuer{
			Value: provider,
		},
		Subject: &saml.Subject{
			NameID: &saml.NameID{
				Value: taxCode,
			},
		},
		Conditions: &saml.Conditions{
			NotBefore:    now,
			NotOnOrAfter: now.Add(5 * time.Minute),
		},
		AttributeStatements: []saml.AttributeStatement{
			{
				Attributes: []saml.Attribute{
					{
						Name: "spidCode",
						Values: []saml.AttributeValue{
							{Value: taxCode},
						},
					},
					{
						Name: "name",
						Values: []saml.AttributeValue{
							{Value: name},
						},
					},
					{
						Name: "familyName",
						Values: []saml.AttributeValue{
							{Value: familyName},
						},
					},
					{
						Name: "email",
						Values: []saml.AttributeValue{
							{Value: email},
						},
					},
				},
			},
		},
	}

	return assertion
}

func TestSPIDAuthenticationFlow(t *testing.T) {
	// Setup
	mockUserRepo := NewMockUserRepository()
	jwtSecret := "test-secret-key"

	// Create mock service provider
	sp := &saml.ServiceProvider{
		EntityID: "https://test.esempio.it",
	}

	spidService := auth.NewSPIDService(mockUserRepo, jwtSecret, sp)

	t.Run("Successful SPID Authentication - New User", func(t *testing.T) {
		// Create mock SAML assertion
		assertion := createMockSAMLAssertion(
			"RSSMRA80A01H501U", // Tax code
			"Mario",            // Name
			"Rossi",            // Family name
			"mario.rossi@example.com",
			"https://posteid.poste.it", // Poste SPID
		)

		// Validate assertion and extract attributes
		attrs, err := spidService.ValidateSAMLAssertion(context.Background(), assertion)
		require.NoError(t, err)

		// Verify attributes
		assert.Equal(t, "RSSMRA80A01H501U", attrs.TaxCode)
		assert.Equal(t, "Mario", attrs.Name)
		assert.Equal(t, "Rossi", attrs.FamilyName)
		assert.Equal(t, "mario.rossi@example.com", attrs.Email)
		assert.Contains(t, attrs.Provider, "poste")

		// Create or get user
		schoolID := uint(1)
		user, err := spidService.GetOrCreateUserBySPID(context.Background(), attrs, schoolID)
		require.NoError(t, err)

		// Verify user was created
		assert.NotNil(t, user)
		assert.Equal(t, "mario.rossi@example.com", user.Email)
		assert.Equal(t, "spid", user.AuthMethod)
		assert.NotNil(t, user.SPIDProvider)
		assert.NotNil(t, user.ExternalID)
		assert.Contains(t, *user.ExternalID, "spid:")
	})

	t.Run("Successful SPID Authentication - Existing User", func(t *testing.T) {
		// Create existing user
		externalID := "spid:BNCGVN85R02H501Z"
		existingUser := &domain.User{
			ID:         1,
			Email:      "giovanni.bianchi@example.com",
			AuthMethod: "spid",
			ExternalID: &externalID,
			SchoolID:   1,
		}
		mockUserRepo.Create(existingUser)

		// Create SAML assertion for same user
		assertion := createMockSAMLAssertion(
			"BNCGVN85R02H501Z",
			"Giovanni",
			"Bianchi",
			"giovanni.bianchi@example.com",
			"https://login.aruba.it",
		)

		attrs, err := spidService.ValidateSAMLAssertion(context.Background(), assertion)
		require.NoError(t, err)

		schoolID := uint(1)
		user, err := spidService.GetOrCreateUserBySPID(context.Background(), attrs, schoolID)
		require.NoError(t, err)

		// Should return existing user
		assert.Equal(t, existingUser.Email, user.Email)
		assert.NotNil(t, user.LastAuthAt)
	})

	t.Run("Invalid SAML Assertion - Missing Tax Code", func(t *testing.T) {
		assertion := createMockSAMLAssertion(
			"", // Missing tax code
			"Test",
			"User",
			"test@example.com",
			"https://test.provider.it",
		)

		_, err := spidService.ValidateSAMLAssertion(context.Background(), assertion)
		assert.Error(t, err)
		assert.Equal(t, auth.ErrSPIDAttributeMissing, err)
	})

	t.Run("Invalid SAML Assertion - Expired", func(t *testing.T) {
		assertion := createMockSAMLAssertion(
			"RSSMRA80A01H501U",
			"Mario",
			"Rossi",
			"mario.rossi@example.com",
			"https://posteid.poste.it",
		)

		// Set expiry to past
		pastTime := time.Now().Add(-10 * time.Minute)
		assertion.Conditions.NotOnOrAfter = pastTime

		_, err := spidService.ValidateSAMLAssertion(context.Background(), assertion)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expired")
	})
}

func TestSPIDProviders(t *testing.T) {
	providers := []string{
		"https://posteid.poste.it",
		"https://login.aruba.it",
		"https://loginspid.infocert.it",
		"https://id.lepida.it",
		"https://login.id.tim.it",
		"https://spid.register.it",
		"https://identity.sieltecloud.it",
		"https://spid.intesa.it",
	}

	for _, provider := range providers {
		t.Run("Provider: "+provider, func(t *testing.T) {
			assertion := createMockSAMLAssertion(
				"RSSMRA80A01H501U",
				"Mario",
				"Rossi",
				"mario.rossi@example.com",
				provider,
			)

			assert.NotNil(t, assertion)
			assert.Equal(t, provider, assertion.Issuer.Value)
		})
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
