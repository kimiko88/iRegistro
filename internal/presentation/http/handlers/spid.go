package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/crewjam/saml"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/k/iRegistro/internal/application/auth"
)

type SPIDHandler struct {
	spidService     *auth.SPIDService
	serviceProvider *saml.ServiceProvider
}

func NewSPIDHandler(spidService *auth.SPIDService, sp *saml.ServiceProvider) *SPIDHandler {
	return &SPIDHandler{
		spidService:     spidService,
		serviceProvider: sp,
	}
}

// Login initiates SPID authentication flow
// GET /auth/spid/login?redirect_uri=<url>&school_id=<uuid>
func (h *SPIDHandler) Login(c *gin.Context) {
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

	// Store school_id in relay state for callback
	relayState := fmt.Sprintf("%s|%s", redirectURI, schoolID)

	// Create SAML authentication request
	binding := saml.HTTPRedirectBinding
	bindingLocation := h.serviceProvider.GetSSOBindingLocation(binding)

	req, err := h.serviceProvider.MakeAuthenticationRequest(bindingLocation, binding, saml.HTTPPostBinding)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SAML request"})
		return
	}

	// Redirect to SPID aggregator
	redirectURL, err := binding.GetSignedRedirectURL(req, relayState)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to sign SAML request"})
		return
	}

	c.Redirect(http.StatusFound, redirectURL.String())
}

// Callback handles SPID SAML assertion callback
// POST /auth/spid/callback
func (h *SPIDHandler) Callback(c *gin.Context) {
	// Parse SAML response
	err := c.Request.ParseForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	// Get relay state (contains redirect_uri and school_id)
	relayState := c.Request.Form.Get("RelayState")

	// Parse SAML response
	assertion, err := h.serviceProvider.ParseResponse(c.Request, []string{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SAML validation failed", "details": err.Error()})
		return
	}

	// Validate and extract SPID attributes
	spidAttrs, err := h.spidService.ValidateSAMLAssertion(c.Request.Context(), assertion)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "SPID attribute extraction failed", "details": err.Error()})
		return
	}

	// Parse relay state to get school_id
	var redirectURI string
	var schoolID uint

	if relayState != "" {
		// Format: "redirect_uri|school_id"
		var schoolIDStr string
		fmt.Sscanf(relayState, "%s|%s", &redirectURI, &schoolIDStr)

		// Parse school_id as uint
		var err error
		var id uint64
		id, err = strconv.ParseUint(schoolIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid school_id in relay state"})
			return
		}
		schoolID = uint(id)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing relay state"})
		return
	}

	// Get or create user
	user, err := h.spidService.GetOrCreateUserBySPID(c.Request.Context(), spidAttrs, schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed", "details": err.Error()})
		return
	}

	// Generate JWT token
	token, err := h.spidService.GenerateTokenForUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	// Redirect to frontend with token
	if redirectURI == "" {
		redirectURI = os.Getenv("FRONTEND_URL") + "/auth-callback"
	}

	finalURL := fmt.Sprintf("%s?token=%s", redirectURI, token)
	c.Redirect(http.StatusFound, finalURL)
}

// Metadata returns SPID service provider metadata
// GET /auth/spid/metadata
func (h *SPIDHandler) Metadata(c *gin.Context) {
	metadata := h.serviceProvider.Metadata()
	c.Data(http.StatusOK, "application/xml", metadata)
}
