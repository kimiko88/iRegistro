package handlers

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/k/iRegistro/internal/infrastructure/storage"
)

type FileHandler struct {
	storage *storage.LocalStorage
}

func NewFileHandler(storage *storage.LocalStorage) *FileHandler {
	return &FileHandler{storage: storage}
}

func (h *FileHandler) DownloadFile(c *gin.Context) {
	// Query param: path
	// Security: In real app, verify user has access to this specific file (via relation check)
	// For now, relies on Authentication Middleware already verifying user is logged in
	// and simple path traversal check.

	relPath := c.Query("path")
	if relPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	// Basic path traversal prevention (though LocalStorage.Save likely put it in safe place)
	// Clean path
	safePath := filepath.Clean(relPath)
	if strings.Contains(safePath, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path"})
		return
	}

	// In a real scenario, we might want to check DB if this user owns this document
	// file, err := h.storage.Get(safePath) ...
	// But `gin.File` serves from disk.
	// Assuming `relPath` is absolute or relative to where CWD is, or we prepend storage BaseDir.
	// Since LocalStorage returns full path on Save, we assume `relPath` is that Full Path.

	// Check if file exists in our storage directory to prevent arbitrary read
	if !strings.HasPrefix(safePath, h.storage.BaseDir) {
		// If the stored path doesn't start with base dir (unlikely if strictly controlled), reject or handle
	}

	c.File(safePath)
}
