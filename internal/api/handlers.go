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

type Metrics interface {
	IncPastesCreated()
}

type API struct {
	storage *storage.Storage
	metrics Metrics
}

func New(s *storage.Storage, m Metrics) *API {
	return &API{storage: s, metrics: m}
}

type ErrorResponse struct {
	Message string `json:"message" example:"error message"`
}

func (a *API) RegisterRoutes(router *gin.Engine) {
	router.GET("/health", a.healthCheck)

	v1 := router.Group("/api/v1")
	{
		v1.POST("/paste", a.createPaste)
		v1.GET("/paste/:id", a.getPaste)
	}
}

func (a *API) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

// @Summary      Create a new paste
// @Description  Saves a text snippet to the database with an expiration time
// @Tags         pastes
// @Accept       json
// @Produce      json
// @Param        paste   body      CreatePasteRequest  true  "Paste Data"
// @Success      201     {object}  CreatePasteResponse
// @Failure      400     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /paste [post]
func (a *API) createPaste(c *gin.Context) {
	var req CreatePasteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}

	a.metrics.IncPastesCreated()

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
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "failed to save paste"})
		return
	}

	c.JSON(http.StatusCreated, CreatePasteResponse{
		ID:  newPaste.ID,
		URL: "/p/" + newPaste.ID,
	})
}

// @Summary      Get a paste by ID
// @Description  Retrieves a text snippet from the database by its short ID
// @Tags         pastes
// @Produce      json
// @Param        id   path      string  true  "Paste ID"
// @Success      200  {object}  paste.Paste
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /paste/{id} [get]
func (a *API) getPaste(c *gin.Context) {
	id := c.Param("id")

	p, err := a.storage.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: "paste not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "database error"})
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
	rand.Read(b)
	for i := 0; i < n; i++ {
		b[i] = letterBytes[int(b[i])%len(letterBytes)]
	}
	return string(b)
}
