package api

import (
	"crypto/rand"
	"errors"
	"my-pastebin/internal/paste"
	"my-pastebin/internal/storage"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreatePasteRequest struct {
	Content   string `json:"content" binding:"required"`
	ExpiresIn string `json:"expires_in"` // "10m", "1h", "1d", "never"
	Syntax    string `json:"syntax"`
}

type CreatePasteResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type API struct {
	storage *storage.Storage
}

func New(s *storage.Storage) *API {
	return &API{storage: s}
}

func (a *API) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		v1.POST("/paste", a.createPaste)
		v1.GET("/paste/:id", a.getPaste)
	}
}

func (a *API) createPaste(c *gin.Context) {
	var req CreatePasteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var expiresAt *time.Time
	if d, err := parseDuration(req.ExpiresIn); err == nil {
		t := time.Now().Add(d)
		expiresAt = &t
	}

	newPaste := &paste.Paste{
		ID:        generateShortID(8),
		Content:   req.Content,
		Syntax:    req.Syntax,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}

	if err := a.storage.Save(c.Request.Context(), newPaste); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save paste"})
		return
	}

	c.JSON(http.StatusCreated, CreatePasteResponse{
		ID:  newPaste.ID,
		URL: "/p/" + newPaste.ID,
	})
}

func (a *API) getPaste(c *gin.Context) {
	id := c.Param("id")

	p, err := a.storage.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "paste not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, p)
}

func parseDuration(s string) (time.Duration, error) {
	switch s {
	case "10m":
		return 10 * time.Minute, nil
	case "1h":
		return time.Hour, nil
	case "1d":
		return 24 * time.Hour, nil
	default:
		return 0, errors.New("invalid duration")
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortID(n int) string {
	b := make([]byte, n)
	rand.Read(b) // Используем крипто-стойкий генератор
	for i := 0; i < n; i++ {
		b[i] = letterBytes[int(b[i])%len(letterBytes)]
	}
	return string(b)
}
