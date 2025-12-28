package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/k/iRegistro/internal/application/auth"
	"github.com/k/iRegistro/internal/domain"
)

type AuthHandler struct {
	service *auth.AuthService
}

func NewAuthHandler(s *auth.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

type RegisterRequest struct {
	Email     string      `json:"email" binding:"required,email"`
	Password  string      `json:"password" binding:"required,min=8"`
	Role      domain.Role `json:"role" binding:"required"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	SchoolID  uint        `json:"school_id"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	OTPCode  string `json:"otp_code"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &domain.User{
		Email:     req.Email,
		Role:      req.Role,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		SchoolID:  req.SchoolID,
	}

	if err := h.service.Register(user, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, refreshToken, err := h.service.Login(req.Email, req.Password, req.OTPCode, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	// Set Refresh Token in HttpOnly Cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true, // Should be true in production
		Path:     "/auth/refresh",
		SameSite: http.SameSiteStrictMode,
	})

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"expires_in":   900, // 15 minutes
	})
}

// 2FA Handlers

func (h *AuthHandler) Enable2FA(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	secret, url, err := h.service.Enable2FA(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to enable 2fa"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"secret": secret,
		"qr_url": url,
	})
}

type Verify2FARequest struct {
	Code string `json:"code" binding:"required"`
}

func (h *AuthHandler) Verify2FA(c *gin.Context) {
	var req Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.service.VerifyAndEnable2FA(userID.(uint), req.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2fa enabled successfully"})
}

// Password Reset Handlers

type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) RequestPasswordReset(c *gin.Context) {
	var req PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.RequestPasswordReset(req.Email); err != nil {
		// Log error but return success to avoid enumeration?
		// For dev, return error maybe.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "if email exists, reset instructions sent"})
}

type ResetPasswordConfirmRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ResetPassword(req.Email, req.Token, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // Could be invalid token
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}
