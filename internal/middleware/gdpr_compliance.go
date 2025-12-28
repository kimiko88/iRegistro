package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// GDPRComplianceMiddleware adds GDPR compliance headers and tracking
func GDPRComplianceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add security headers for GDPR compliance

		// HSTS (HTTP Strict Transport Security) - forces HTTPS
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// X-Content-Type-Options prevents MIME sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// X-Frame-Options prevents clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Content-Security-Policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'")

		// X-XSS-Protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer-Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions-Policy (formerly Feature-Policy)
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Store IP and User-Agent for GDPR audit logging
		clientIP := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")

		// Attach to context for use in GDPR services
		c.Set("ip_address", clientIP)
		c.Set("user_agent", userAgent)

		c.Next()
	}
}

// TLSRequiredMiddleware ensures TLS 1.3 minimum in production
func TLSRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// In production, reject non-HTTPS requests
		if c.Request.TLS == nil && gin.Mode() == gin.ReleaseMode {
			c.AbortWithStatusJSON(426, gin.H{
				"error": "Upgrade Required: HTTPS only",
			})
			return
		}

		// Check TLS version (1.3 minimum)
		if c.Request.TLS != nil && c.Request.TLS.Version < 0x0304 { // TLS 1.3
			c.AbortWithStatusJSON(426, gin.H{
				"error": "TLS 1.3 required",
			})
			return
		}

		c.Next()
	}
}

// DataMinimizationMiddleware logs potentially excessive data requests
func DataMinimizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for overly broad queries
		query := c.Request.URL.Query()

		// Warn if requesting all users without pagination
		if strings.Contains(c.Request.URL.Path, "/users") && !query.Has("limit") {
			c.Header("X-GDPR-Warning", "Data minimization: Consider using pagination")
		}

		c.Next()
	}
}
