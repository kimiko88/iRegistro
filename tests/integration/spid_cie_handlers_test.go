package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSPIDHTTPHandlerFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("SPID Login Redirect", func(t *testing.T) {
		router := gin.New()

		// Mock SPID handler (simplified for testing)
		router.GET("/auth/spid/login", func(c *gin.Context) {
			schoolID := c.Query("school_id")
			redirectURI := c.Query("redirect_uri")

			if schoolID == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "school_id required"})
				return
			}

			// In real implementation, this would create SAML request
			// For test, just verify parameters
			c.JSON(http.StatusOK, gin.H{
				"school_id":    schoolID,
				"redirect_uri": redirectURI,
				"status":       "redirect_to_spid",
			})
		})

		req := httptest.NewRequest("GET", "/auth/spid/login?school_id=test-uuid&redirect_uri=http://localhost:3000/callback", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "redirect_to_spid")
	})

	t.Run("SPID Login Missing School ID", func(t *testing.T) {
		router := gin.New()

		router.GET("/auth/spid/login", func(c *gin.Context) {
			schoolID := c.Query("school_id")

			if schoolID == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "school_id required"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/auth/spid/login", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "school_id required")
	})

	t.Run("SPID Metadata Endpoint", func(t *testing.T) {
		router := gin.New()

		router.GET("/auth/spid/metadata", func(c *gin.Context) {
			// Return mock XML metadata
			metadata := `<?xml version="1.0"?>
<EntityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" entityID="https://test.esempio.it">
  <SPSSODescriptor>
    <AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://test.esempio.it/auth/spid/callback"/>
  </SPSSODescriptor>
</EntityDescriptor>`

			c.Data(http.StatusOK, "application/xml", []byte(metadata))
		})

		req := httptest.NewRequest("GET", "/auth/spid/metadata", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/xml", w.Header().Get("Content-Type"))
		assert.Contains(t, w.Body.String(), "EntityDescriptor")
		assert.Contains(t, w.Body.String(), "AssertionConsumerService")
	})
}

func TestCIEHTTPHandlerFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("CIE Login Redirect", func(t *testing.T) {
		router := gin.New()

		router.GET("/auth/cie/login", func(c *gin.Context) {
			schoolID := c.Query("school_id")
			redirectURI := c.Query("redirect_uri")

			if schoolID == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "school_id required"})
				return
			}

			// Set state cookie
			c.SetCookie("oauth_state", "test_state|"+redirectURI+"|"+schoolID, 300, "/", "", false, true)

			c.JSON(http.StatusOK, gin.H{
				"auth_url": "https://mock-cie-provider.it/authorize?state=test_state",
			})
		})

		req := httptest.NewRequest("GET", "/auth/cie/login?school_id=test-uuid&redirect_uri=http://localhost:3000/callback", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Check cookie was set
		cookies := w.Result().Cookies()
		var stateCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "oauth_state" {
				stateCookie = cookie
				break
			}
		}

		require.NotNil(t, stateCookie)
		assert.Contains(t, stateCookie.Value, "test_state")
		assert.True(t, stateCookie.HttpOnly)
	})

	t.Run("CIE Callback State Validation", func(t *testing.T) {
		router := gin.New()

		router.GET("/auth/cie/callback", func(c *gin.Context) {
			// Get state from cookie
			stateCookie, err := c.Cookie("oauth_state")
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Missing state cookie"})
				return
			}

			// Get state from query
			stateQuery := c.Query("state")
			if stateQuery == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Missing state parameter"})
				return
			}

			// Verify states match (simplified)
			if stateCookie != stateQuery {
				c.JSON(http.StatusBadRequest, gin.H{"error": "State mismatch"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "authenticated"})
		})

		req := httptest.NewRequest("GET", "/auth/cie/callback?state=valid_state&code=auth_code", nil)
		req.AddCookie(&http.Cookie{
			Name:  "oauth_state",
			Value: "valid_state",
		})
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "authenticated")
	})

	t.Run("CIE Callback CSRF Protection - State Mismatch", func(t *testing.T) {
		router := gin.New()

		router.GET("/auth/cie/callback", func(c *gin.Context) {
			stateCookie, _ := c.Cookie("oauth_state")
			stateQuery := c.Query("state")

			if stateCookie != stateQuery {
				c.JSON(http.StatusBadRequest, gin.H{"error": "State mismatch - possible CSRF attack"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/auth/cie/callback?state=malicious_state&code=auth_code", nil)
		req.AddCookie(&http.Cookie{
			Name:  "oauth_state",
			Value: "valid_state",
		})
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "CSRF")
	})
}

func TestAuthenticationSecurityHeaders(t *testing.T) {
	t.Run("HTTPS Redirect in Production", func(t *testing.T) {
		// In production, non-HTTPS requests should be rejected
		// This is typically handled by middleware or reverse proxy

		router := gin.New()
		router.Use(func(c *gin.Context) {
			if c.Request.TLS == nil && gin.Mode() == gin.ReleaseMode {
				c.AbortWithStatusJSON(http.StatusUpgradeRequired, gin.H{
					"error": "HTTPS required",
				})
				return
			}
			c.Next()
		})

		router.GET("/auth/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		gin.SetMode(gin.ReleaseMode)

		req := httptest.NewRequest("GET", "/auth/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Should require HTTPS in production
		assert.Equal(t, http.StatusUpgradeRequired, w.Code)

		gin.SetMode(gin.TestMode) // Reset
	})
}
