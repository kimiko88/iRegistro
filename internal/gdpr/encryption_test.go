package gdpr_test

import (
	"testing"

	"github.com/k/iRegistro/internal/gdpr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncryption(t *testing.T) {
	// Set test encryption key (32 bytes for AES-256)
	t.Setenv("ENCRYPTION_KEY", "12345678901234567890123456789012")

	encSvc, err := gdpr.NewEncryptionService()
	require.NoError(t, err)

	t.Run("Encrypt and Decrypt", func(t *testing.T) {
		plaintext := "RSSMRA80A01H501U" // Example Italian tax code

		encrypted, err := encSvc.Encrypt(plaintext)
		require.NoError(t, err)
		assert.NotEmpty(t, encrypted)
		assert.NotEqual(t, plaintext, encrypted)

		decrypted, err := encSvc.Decrypt(encrypted)
		require.NoError(t, err)
		assert.Equal(t, plaintext, decrypted)
	})

	t.Run("EncryptTaxCode", func(t *testing.T) {
		taxCode := "RSSMRA80A01H501U"

		encrypted, err := encSvc.EncryptTaxCode(taxCode)
		require.NoError(t, err)

		decrypted, err := encSvc.DecryptTaxCode(encrypted)
		require.NoError(t, err)
		assert.Equal(t, taxCode, decrypted)
	})

	t.Run("EncryptPhone", func(t *testing.T) {
		phone := "+39 333 1234567"

		encrypted, err := encSvc.EncryptPhone(phone)
		require.NoError(t, err)

		decrypted, err := encSvc.DecryptPhone(encrypted)
		require.NoError(t, err)
		assert.Equal(t, phone, decrypted)
	})

	t.Run("EncryptAddress", func(t *testing.T) {
		address := "Via Roma 123, 00100 Roma RM"

		encrypted, err := encSvc.EncryptAddress(address)
		require.NoError(t, err)

		decrypted, err := encSvc.DecryptAddress(encrypted)
		require.NoError(t, err)
		assert.Equal(t, address, decrypted)
	})

	t.Run("Empty String", func(t *testing.T) {
		encrypted, err := encSvc.Encrypt("")
		require.NoError(t, err)
		assert.Empty(t, encrypted)

		decrypted, err := encSvc.Decrypt("")
		require.NoError(t, err)
		assert.Empty(t, decrypted)
	})
}

func TestInvalidKey(t *testing.T) {
	t.Run("Key Too Short", func(t *testing.T) {
		t.Setenv("ENCRYPTION_KEY", "short")

		_, err := gdpr.NewEncryptionService()
		assert.Error(t, err)
		assert.Equal(t, gdpr.ErrInvalidKey, err)
	})

	t.Run("Missing Key", func(t *testing.T) {
		t.Setenv("ENCRYPTION_KEY", "")

		_, err := gdpr.NewEncryptionService()
		assert.Error(t, err)
	})
}
