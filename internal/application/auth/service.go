package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/k/iRegistro/internal/domain"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	UserID   uint        `json:"user_id"`
	SchoolID uint        `json:"school_id"`
	Role     domain.Role `json:"role"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // Cost 14 is good for security
	return string(bytes), err
}

func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GenerateAccessToken(user *domain.User, secret string, duration time.Duration) (string, *CustomClaims, error) {
	expirationTime := time.Now().Add(duration)
	claims := &CustomClaims{
		UserID:   user.ID,
		SchoolID: user.SchoolID,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", nil, err
	}

	return tokenString, claims, nil
}

type AuthService struct {
	userRepo        domain.UserRepository
	authRepo        domain.AuthRepository
	jwtSecret       string
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func NewAuthService(u domain.UserRepository, a domain.AuthRepository, secret string, accessDur, refreshDur time.Duration) *AuthService {
	return &AuthService{
		userRepo:        u,
		authRepo:        a,
		jwtSecret:       secret,
		accessDuration:  accessDur,
		refreshDuration: refreshDur,
	}
}

func (s *AuthService) Register(user *domain.User, password string) error {
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}
	user.PasswordHash = hash
	return s.userRepo.Create(user)
}

func (s *AuthService) Login(email, password, otpCode, ip, userAgent string) (*domain.User, string, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", "", err
	}
	if user == nil {
		return nil, "", "", domain.ErrInvalidCredentials
	}

	// Check Lockout
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		return nil, "", "", fmt.Errorf("account locked until %s", user.LockedUntil.Format(time.RFC3339))
	}

	if err := CheckPassword(password, user.PasswordHash); err != nil {
		// Increment Failed Attempts
		user.FailedLogins++
		if user.FailedLogins >= 5 {
			lockTime := time.Now().Add(15 * time.Minute)
			user.LockedUntil = &lockTime
		}
		s.userRepo.Update(user)
		return nil, "", "", domain.ErrInvalidCredentials
	}

	// Check 2FA
	if user.TwoFAEnabled {
		if otpCode == "" {
			return nil, "", "", fmt.Errorf("2fa required") // Should define proper error
		}
		if !totp.Validate(otpCode, user.TwoFASecret) {
			return nil, "", "", fmt.Errorf("invalid 2fa code")
		}
	}

	// Reset Failed Attempts on success
	if user.FailedLogins > 0 || user.LockedUntil != nil {
		user.FailedLogins = 0
		user.LockedUntil = nil
		s.userRepo.Update(user)
	}

	// Generate Access Token
	accessToken, _, err := GenerateAccessToken(user, s.jwtSecret, s.accessDuration)
	if err != nil {
		return nil, "", "", err
	}

	// Generate Refresh Token (simplified for now, usually a random string)
	refreshToken, _, err := GenerateAccessToken(user, s.jwtSecret, s.refreshDuration)
	// Ideally Refresh Token is opaque, but for simplicity reusing JWT or random hex
	if err != nil {
		return nil, "", "", err
	}

	// Store Session
	session := &domain.Session{
		UserID:    user.ID,
		TokenHash: accessToken, // Or hash of it
		ExpiresAt: time.Now().Add(s.accessDuration),
		IPAddress: ip,
		UserAgent: userAgent,
	}
	if err := s.authRepo.CreateSession(session); err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

// 2FA Methods

func (s *AuthService) Enable2FA(userID uint) (string, string, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", domain.ErrUserNotFound
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "iRegistro",
		AccountName: user.Email,
	})
	if err != nil {
		return "", "", err
	}

	// Store secret temporarily or permanently?
	// Usually we store it but mark 2FA as disabled until verified.
	// But to avoid "junk" secrets, we might return it and only save on verify.
	// However, to verify we need the secret.
	// We'll save it now but keep TwoFAEnabled = false
	user.TwoFASecret = key.Secret()
	if err := s.userRepo.Update(user); err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

func (s *AuthService) VerifyAndEnable2FA(userID uint, code string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return domain.ErrUserNotFound
	}

	if !totp.Validate(code, user.TwoFASecret) {
		return fmt.Errorf("invalid OTP code")
	}

	user.TwoFAEnabled = true
	return s.userRepo.Update(user)
}

func (s *AuthService) Verify2FALogin(userID uint, code string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return domain.ErrUserNotFound
	}

	if !user.TwoFAEnabled {
		return nil // Or error? Logic depends on flow.
	}

	if !totp.Validate(code, user.TwoFASecret) {
		return fmt.Errorf("invalid OTP code")
	}
	return nil
}

// Password Reset

func (s *AuthService) RequestPasswordReset(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return nil // Return nil to avoid email enumeration
	}

	// Generate Reset Token (Random string)
	token := "mock-reset-token-" + email // In prod use crypto/rand

	// Hash token
	hash, err := HashPassword(token) // Reuse bcrypt helper
	if err != nil {
		return err
	}

	user.ResetTokenHash = hash
	expiration := time.Now().Add(15 * time.Minute)
	user.ResetTokenExp = &expiration

	if err := s.userRepo.Update(user); err != nil {
		return err
	}

	// Send Email (Mock)
	fmt.Printf("SEND EMAIL TO %s: Reset Token: %s\n", email, token)
	return nil
}

func (s *AuthService) ResetPassword(email, token, newPassword string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return domain.ErrUserNotFound
	}

	if user.ResetTokenExp == nil || time.Now().After(*user.ResetTokenExp) {
		return fmt.Errorf("token expired")
	}

	if err := CheckPassword(token, user.ResetTokenHash); err != nil {
		return fmt.Errorf("invalid token")
	}

	newHash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = newHash
	user.ResetTokenHash = ""
	user.ResetTokenExp = nil

	return s.userRepo.Update(user)
}

func (s *AuthService) GetUserByID(id uint) (*domain.User, error) {
	return s.userRepo.FindByID(id)
}
