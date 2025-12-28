package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k/iRegistro/internal/application/auth"
)

type CIEHandler struct {
	cieService *auth.CIEService
}

func NewCIEHandler(cieService *auth.CIEService) *CIEHandler {
	return &CIEHandler{
		cieService: cieService,
	}
}

// Login initiates CIE authentication flow
// GET /auth/cie/login?redirect_uri=<url>&school_id=<uuid>
func (h *CIEHandler) Login(c *gin.Context) {
	redirectURI := c.Query("redirect_uri")
	if redirectURI == "" {
		redirectURI = os.Getenv("FRONTEND_URL") + "/auth-callback"
	}

	schoolID := c.Query("school_id")
	if schoolID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "school_id required"})
		return
	}

	// Validate school_id is UUID
	if _, err := uuid.Parse(schoolID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid school_id"})
		return
	}

	// Generate and store state for CSRF protection
	state := generateRandomState()

	// Store state in secure cookie with redirect_uri and school_id
	stateData := fmt.Sprintf("%s|%s|%s", state, redirectURI, schoolID)
	c.SetCookie(
		"oauth_state",
		stateData,
		300, // 5 minutes
		"/",
		"",
		true, // Secure (HTTPS only)
		true, // HttpOnly
	)

	// Get OAuth2 authorization URL
	authURL := h.cieService.GetAuthorizationURL(state)

	c.Redirect(http.StatusFound, authURL)
}

// Callback handles CIE OIDC callback
// GET /auth/cie/callback?code=<code>&state=<state>
func (h *CIEHandler) Callback(c *gin.Context) {
	// Verify state to prevent CSRF
	stateCookie, err := c.Cookie("oauth_state")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing state cookie"})
		return
	}

	// Clear state cookie
	c.SetCookie("oauth_state", "", -1, "/", "", true, true)

	stateQuery := c.Query("state")
	if stateQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing state parameter"})
		return
	}

	// Parse state cookie: "state|redirect_uri|school_id"
	var state, redirectURI, schoolIDStr string
	fmt.Sscanf(stateCookie, "%s|%s|%s", &state, &redirectURI, &schoolIDStr)

	// Verify state matches
	if state != stateQuery {
		c.JSON(http.StatusBadRequest, gin.H{"error": "State mismatch - possible CSRF attack"})
		return
	}

	// Parse school_id as uint
	id, err := strconv.ParseUint(schoolIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid school_id"})
		return
	}
	schoolID := uint(id)

	// Get authorization code
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing authorization code"})
		return
	}

	// Exchange code for token
	token, err := h.cieService.ExchangeCodeForToken(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token exchange failed", "details": err.Error()})
		return
	}

	// Extract and validate ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No ID token in response"})
		return
	}

	// Validate and extract CIE claims
	claims, err := h.cieService.ValidateAndExtractClaims(c.Request.Context(), rawIDToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token validation failed", "details": err.Error()})
		return
	}

	// Get or create user
	user, err := h.cieService.GetOrCreateUserByCIE(c.Request.Context(), claims, schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed", "details": err.Error()})
		return
	}

	// Generate JWT token
	jwtToken, err := h.cieService.GenerateTokenForUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "JWT generation failed"})
		return
	}

	// Redirect to frontend with JWT token
	if redirectURI == "" {
		redirectURI = os.Getenv("FRONTEND_URL") + "/auth-callback"
	}

	finalURL := fmt.Sprintf("%s?token=%s", redirectURI, jwtToken)
	c.Redirect(http.StatusFound, finalURL)
}
