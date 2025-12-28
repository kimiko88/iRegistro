package auth

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/k/iRegistro/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestPasswordHashing(t *testing.T) {
	password := "secret123"

	// Test Hash
	hash, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	// Test Check Correct
	err = CheckPassword(password, hash)
	assert.NoError(t, err)

	// Test Check Incorrect
	err = CheckPassword("wrongpassword", hash)
	assert.Error(t, err)
}

func TestTokenGeneration(t *testing.T) {
	secret := "testsecret"
	user := &domain.User{
		ID:       1,
		Email:    "test@example.com",
		Role:     domain.RoleTeacher,
		SchoolID: 10,
	}

	// Generate Token
	tokenString, claims, err := GenerateAccessToken(user, secret, 15*time.Minute)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)
	assert.Equal(t, uint(1), claims.UserID)
	assert.Equal(t, uint(10), claims.SchoolID)
	assert.Equal(t, domain.RoleTeacher, claims.Role)

	// Verify Token
	parsedToken, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	parsedClaims, ok := parsedToken.Claims.(*CustomClaims)
	assert.True(t, ok)
	assert.Equal(t, user.ID, parsedClaims.UserID)
}

// Mocks
type MockUserRepository struct {
	users map[string]*domain.User
	err   error
}

func (m *MockUserRepository) Create(user *domain.User) error {
	if m.err != nil {
		return m.err
	}
	m.users[user.Email] = user
	return nil
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.users[email], nil
}

func (m *MockUserRepository) FindByID(id uint) (*domain.User, error) { return nil, nil }
func (m *MockUserRepository) Update(user *domain.User) error         { return nil }
func (m *MockUserRepository) GetByExternalID(ctx context.Context, externalID string) (*domain.User, error) {
	return nil, nil
}
func (m *MockUserRepository) FindAll(schoolID uint) ([]domain.User, error) { return nil, nil }

type MockAuthRepository struct {
	sessions []*domain.Session
}

func (m *MockAuthRepository) CreateSession(session *domain.Session) error {
	m.sessions = append(m.sessions, session)
	return nil
}

func (m *MockAuthRepository) StoreRefreshToken(token *domain.RefreshToken) error { return nil }
func (m *MockAuthRepository) RevokeRefreshToken(tokenHash string) error          { return nil }
func (m *MockAuthRepository) GetRefreshToken(tokenHash string) (*domain.RefreshToken, error) {
	return nil, nil
}

// Service Tests
func TestRegister(t *testing.T) {
	mockUserRepo := &MockUserRepository{users: make(map[string]*domain.User)}
	mockAuthRepo := &MockAuthRepository{}
	service := NewAuthService(mockUserRepo, mockAuthRepo, "secret", 15*time.Minute, 7*24*time.Hour)

	user := &domain.User{
		Email:        "new@example.com",
		PasswordHash: "password", // In real input this is raw password
		Role:         domain.RoleTeacher,
		SchoolID:     1,
	}

	err := service.Register(user, "password")
	assert.NoError(t, err)

	// Verify hashed password
	savedUser := mockUserRepo.users["new@example.com"]
	assert.NotEqual(t, "password", savedUser.PasswordHash)
}

func TestLogin(t *testing.T) {
	mockUserRepo := &MockUserRepository{users: make(map[string]*domain.User)}
	mockAuthRepo := &MockAuthRepository{}
	service := NewAuthService(mockUserRepo, mockAuthRepo, "secret", 15*time.Minute, 7*24*time.Hour)

	// Setup User
	hash, _ := HashPassword("password")
	user := &domain.User{
		ID:           1,
		Email:        "login@example.com",
		PasswordHash: hash,
		Role:         domain.RoleAdmin,
		SchoolID:     1,
	}
	mockUserRepo.users[user.Email] = user

	// Test Success
	userWrapper, token, refToken, err := service.Login("login@example.com", "password", "", "127.0.0.1", "test-agent")
	assert.NoError(t, err)
	assert.NotNil(t, userWrapper)
	assert.NotEmpty(t, token)
	assert.NotEmpty(t, refToken)

	// Test Wrong Password
	_, _, _, err = service.Login("login@example.com", "wrong", "", "127.0.0.1", "test-agent")
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())

	// Test User Not Found
	_, _, _, err = service.Login("unknown@example.com", "password", "", "127.0.0.1", "test-agent")
	assert.Error(t, err)
}

func TestAccountLockout(t *testing.T) {
	mockUserRepo := &MockUserRepository{users: make(map[string]*domain.User)}
	mockAuthRepo := &MockAuthRepository{}
	service := NewAuthService(mockUserRepo, mockAuthRepo, "secret", 15*time.Minute, 7*24*time.Hour)

	user := &domain.User{
		ID:           1,
		Email:        "lockout@example.com",
		PasswordHash: "hash", // CheckPassword compares with this
		Role:         domain.RoleAdmin,
		SchoolID:     1,
	}
	// Pre-hash for CheckPassword to fail
	hash, _ := HashPassword("correct_password")
	user.PasswordHash = hash

	mockUserRepo.users[user.Email] = user

	// Fail 5 times
	for i := 0; i < 5; i++ {
		_, _, _, err := service.Login("lockout@example.com", "wrong", "", "127.0.0.1", "test-agent")
		assert.Error(t, err)
	}

	// 6th attempt should be locked
	_, _, _, err := service.Login("lockout@example.com", "wrong", "", "127.0.0.1", "test-agent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "account locked")
}

func TestPasswordReset(t *testing.T) {
	mockUserRepo := &MockUserRepository{users: make(map[string]*domain.User)}
	mockAuthRepo := &MockAuthRepository{}
	service := NewAuthService(mockUserRepo, mockAuthRepo, "secret", 15*time.Minute, 7*24*time.Hour)

	user := &domain.User{
		ID:       1,
		Email:    "reset@example.com",
		Role:     domain.RoleAdmin,
		SchoolID: 1,
	}
	// Initial password
	hash, _ := HashPassword("old_password")
	user.PasswordHash = hash
	mockUserRepo.users[user.Email] = user

	// 1. Request Reset
	err := service.RequestPasswordReset("reset@example.com")
	assert.NoError(t, err)

	updatedUser := mockUserRepo.users["reset@example.com"]
	assert.NotEmpty(t, updatedUser.ResetTokenHash)
	assert.NotNil(t, updatedUser.ResetTokenExp)

	// 2. Reset with wrong token
	err = service.ResetPassword("reset@example.com", "wrong_token", "new_password")
	assert.Error(t, err)

	// 3. Reset with correct token
	// We mocked the token generation in service "mock-reset-token-" + email
	correctToken := "mock-reset-token-reset@example.com"
	err = service.ResetPassword("reset@example.com", correctToken, "new_password")
	assert.NoError(t, err)

	// 4. Verify new password
	updatedUser = mockUserRepo.users["reset@example.com"]
	err = CheckPassword("new_password", updatedUser.PasswordHash)
	assert.NoError(t, err)
	assert.Empty(t, updatedUser.ResetTokenHash)
	assert.Nil(t, updatedUser.ResetTokenExp)
}
