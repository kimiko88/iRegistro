package gdpr

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

var (
	ErrInvalidKey        = errors.New("encryption key must be 32 bytes")
	ErrInvalidCiphertext = errors.New("invalid ciphertext")
)

// EncryptionService handles field-level encryption/decryption
type EncryptionService struct {
	key []byte
}

// NewEncryptionService creates a new encryption service
// Key should be 32 bytes for AES-256
func NewEncryptionService() (*EncryptionService, error) {
	key := os.Getenv("ENCRYPTION_KEY")
	if key == "" {
		return nil, errors.New("ENCRYPTION_KEY environment variable not set")
	}

	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		return nil, ErrInvalidKey
	}

	return &EncryptionService{key: keyBytes}, nil
}

// Encrypt encrypts plaintext using AES-256-GCM
func (s *EncryptionService) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES-256-GCM
func (s *EncryptionService) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", ErrInvalidCiphertext
	}

	nonce, cipherdata := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, cipherdata, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// EncryptedField represents a sensitive database field
type EncryptedField struct {
	Plaintext  string
	Ciphertext string
}

// MarshalJSON for encrypted field (returns encrypted value)
func (e EncryptedField) MarshalJSON() ([]byte, error) {
	return []byte(`"***REDACTED***"`), nil
}

// Helper functions for specific field types

// EncryptTaxCode encrypts Italian tax code (Codice Fiscale)
func (s *EncryptionService) EncryptTaxCode(taxCode string) (string, error) {
	return s.Encrypt(taxCode)
}

// DecryptTaxCode decrypts Italian tax code
func (s *EncryptionService) DecryptTaxCode(encrypted string) (string, error) {
	return s.Decrypt(encrypted)
}

// EncryptPhone encrypts phone number
func (s *EncryptionService) EncryptPhone(phone string) (string, error) {
	return s.Encrypt(phone)
}

// DecryptPhone decrypts phone number
func (s *EncryptionService) DecryptPhone(encrypted string) (string, error) {
	return s.Decrypt(encrypted)
}

// EncryptAddress encrypts home address
func (s *EncryptionService) EncryptAddress(address string) (string, error) {
	return s.Encrypt(address)
}

// DecryptAddress decrypts home address
func (s *EncryptionService) DecryptAddress(encrypted string) (string, error) {
	return s.Decrypt(encrypted)
}
