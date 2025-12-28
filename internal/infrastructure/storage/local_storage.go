package storage

import (
	"os"
	"path/filepath"
)

type FileStorage interface {
	Save(filename string, data []byte) (string, error)
	Get(filename string) ([]byte, error)
}

type LocalStorage struct {
	BaseDir string
}

func NewLocalStorage(baseDir string) (*LocalStorage, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{BaseDir: baseDir}, nil
}

func (s *LocalStorage) Save(filename string, data []byte) (string, error) {
	// secure filename or append timestamp to avoid collision
	// For simplicity, we assume caller manages logic or we prepend timestamp
	// Making it unique
	safeName := filepath.Clean(filename)
	// e.g. /uploads/report_123.pdf
	fullPath := filepath.Join(s.BaseDir, safeName)

	// Ensure dir exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", err
	}

	return fullPath, nil
}

func (s *LocalStorage) Get(filepathStr string) ([]byte, error) {
	// simplistic implementation
	return os.ReadFile(filepathStr)
}
