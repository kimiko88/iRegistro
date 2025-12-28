package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/k/iRegistro/internal/infrastructure/storage"
	"github.com/stretchr/testify/assert"
)

func TestDownloadFile_BadRequest_NoPath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	localStorage, _ := storage.NewLocalStorage("/tmp")
	h := NewFileHandler(localStorage)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/files/download", nil)

	h.DownloadFile(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDownloadFile_BadRequest_PathTraversal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	localStorage, _ := storage.NewLocalStorage("/tmp")
	h := NewFileHandler(localStorage)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/files/download?path=../etc/passwd", nil)

	h.DownloadFile(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
