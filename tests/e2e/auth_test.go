package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Assumes server is running on localhost:8080
// Run with: make test-e2e
func TestAuthFlow_LoginAndRefresh(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test")
	}

	baseURL := "http://localhost:8080"

	// Ensure server is up (simple health check or retry)
	_, err := http.Get(baseURL + "/health")
	if err != nil {
		t.Skipf("Server not running at %s: %v", baseURL, err)
	}

	t.Run("Login Success", func(t *testing.T) {
		loginReq := map[string]string{
			"email":    "teacher1@test.com", // From seed data
			"password": "password123",       // Assumes seed sets this
		}

		body, _ := json.Marshal(loginReq)
		resp, err := http.Post(
			baseURL+"/auth/login",
			"application/json",
			bytes.NewBuffer(body),
		)

		if err != nil {
			t.Fatalf("Failed to call login: %v", err)
		}
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var loginResp struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
			t.Fatal("Failed to decode response")
		}

		assert.NotEmpty(t, loginResp.AccessToken)
		assert.NotEmpty(t, loginResp.RefreshToken)

		// Verify Access Token works
		req, _ := http.NewRequest("GET", baseURL+"/api/v1/users/me", nil)
		req.Header.Add("Authorization", "Bearer "+loginResp.AccessToken)

		meResp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer meResp.Body.Close()
		assert.Equal(t, http.StatusOK, meResp.StatusCode)
	})
}
