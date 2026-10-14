package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/matin-sp/-online-library-management-system/internal/model"
	"github.com/matin-sp/-online-library-management-system/internal/service"
)

// AuthorHandler handles HTTP requests related to authors.
type AuthorHandler struct {
	service *service.AuthorService
}

func NewAuthorHandler(service *service.AuthorService) *AuthorHandler {
	return &AuthorHandler{service: service}
}

// authorRequest is the DTO for create/update bodies. The id is never taken
// from the body - it comes from the database (create) or the URL (update).
type authorRequest struct {
	Name    string `json:"name" binding:"required,max=100"`
	Country string `json:"country" binding:"max=50"`
}

const invalidAuthorBody = "invalid request body: name is required (max 100 characters), country is optional (max 50 characters)"

// authorIDFromParam reads and validates the :id URL parameter.
// It writes the 400 response itself, so callers just return when ok is false.
func authorIDFromParam(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return 0, false
	}
	return id, true
}

// respondAuthorError maps service errors to HTTP status codes.
// Unknown errors are logged but not shown to the client (they may contain
// internal details such as SQL text).
func respondAuthorError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrAuthorNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "author not found"})
	case errors.Is(err, service.ErrAuthorHasBooks):
		c.JSON(http.StatusConflict, gin.H{"error": "cannot delete an author who still has books"})
	default:
		log.Println("author handler error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

// Create handles POST /api/authors
func (h *AuthorHandler) Create(c *gin.Context) {
	var req authorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidAuthorBody})
		return
	}

	author, err := h.service.Create(model.Author{Name: req.Name, Country: req.Country})
	if err != nil {
		respondAuthorError(c, err)
		return
	}
	c.JSON(http.StatusCreated, author)
}

// GetAll handles GET /api/authors
func (h *AuthorHandler) GetAll(c *gin.Context) {
	authors, err := h.service.GetAll()
	if err != nil {
		respondAuthorError(c, err)
		return
	}
	c.JSON(http.StatusOK, authors)
}

// GetByID handles GET /api/authors/:id
func (h *AuthorHandler) GetByID(c *gin.Context) {
	id, ok := authorIDFromParam(c)
	if !ok {
		return
	}

	author, err := h.service.GetByID(id)
	if err != nil {
		respondAuthorError(c, err)
		return
	}
	c.JSON(http.StatusOK, author)
}

// Update handles PUT /api/authors/:id (full replacement of name and country)
func (h *AuthorHandler) Update(c *gin.Context) {
	id, ok := authorIDFromParam(c)
	if !ok {
		return
	}

	var req authorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidAuthorBody})
		return
	}

	author := model.Author{ID: id, Name: req.Name, Country: req.Country}
	if err := h.service.Update(author); err != nil {
		respondAuthorError(c, err)
		return
	}
	c.JSON(http.StatusOK, author)
}

// Delete handles DELETE /api/authors/:id
func (h *AuthorHandler) Delete(c *gin.Context) {
	id, ok := authorIDFromParam(c)
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		respondAuthorError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "author deleted successfully"})
}