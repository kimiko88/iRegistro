package handlers

import (
	"crypto/rand"
	"encoding/base64"
)

// generateRandomState generates a cryptographically secure random state for OAuth2/SAML
func generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
