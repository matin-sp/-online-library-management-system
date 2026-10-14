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

// BookHandler handles HTTP requests related to books.
type BookHandler struct {
	service *service.BookService
}

func NewBookHandler(service *service.BookService) *BookHandler {
	return &BookHandler{service: service}
}

// bookRequest is the DTO for create/update bodies.
type bookRequest struct {
	Title       string `json:"title" binding:"required,max=255"`
	ReleaseYear int    `json:"release_year" binding:"gte=0,lte=2100"`
	AuthorID    int    `json:"author_id" binding:"required,gte=1"`
}

const invalidBookBody = "invalid request body: title is required (max 255 characters), author_id is required (positive integer), release_year must be between 0 and 2100"

// bookIDFromParam reads and validates the :id URL parameter.
func bookIDFromParam(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return 0, false
	}
	return id, true
}

// respondBookError maps service errors to HTTP status codes.
func respondBookError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrBookNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
	case errors.Is(err, service.ErrInvalidAuthor):
		c.JSON(http.StatusBadRequest, gin.H{"error": "author_id does not refer to an existing author"})
	default:
		log.Println("book handler error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

// Create handles POST /api/books
func (h *BookHandler) Create(c *gin.Context) {
	var req bookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidBookBody})
		return
	}

	book, err := h.service.Create(model.Book{
		Title:       req.Title,
		ReleaseYear: req.ReleaseYear,
		AuthorID:    req.AuthorID,
	})
	if err != nil {
		respondBookError(c, err)
		return
	}
	c.JSON(http.StatusCreated, book)
}

// GetAll handles GET /api/books and GET /api/books?author_id=3
func (h *BookHandler) GetAll(c *gin.Context) {
	if raw := c.Query("author_id"); raw != "" {
		authorID, err := strconv.Atoi(raw)
		if err != nil || authorID < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "author_id must be a positive integer"})
			return
		}
		books, err := h.service.GetByAuthorID(authorID)
		if err != nil {
			respondBookError(c, err)
			return
		}
		c.JSON(http.StatusOK, books)
		return
	}

	books, err := h.service.GetAll()
	if err != nil {
		respondBookError(c, err)
		return
	}
	c.JSON(http.StatusOK, books)
}

// GetByID handles GET /api/books/:id
func (h *BookHandler) GetByID(c *gin.Context) {
	id, ok := bookIDFromParam(c)
	if !ok {
		return
	}

	book, err := h.service.GetByID(id)
	if err != nil {
		respondBookError(c, err)
		return
	}
	c.JSON(http.StatusOK, book)
}

// GetByIDWithAuthor handles GET /api/books/:id/with-author
// (one JOIN query returning the book plus its author's name)
func (h *BookHandler) GetByIDWithAuthor(c *gin.Context) {
	id, ok := bookIDFromParam(c)
	if !ok {
		return
	}

	result, err := h.service.GetByIDWithAuthor(id)
	if err != nil {
		respondBookError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"book":        result.Book,
		"author_name": result.AuthorName,
	})
}

// Update handles PUT /api/books/:id (full replacement of the book's fields)
func (h *BookHandler) Update(c *gin.Context) {
	id, ok := bookIDFromParam(c)
	if !ok {
		return
	}

	var req bookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidBookBody})
		return
	}

	book := model.Book{
		ID:          id,
		Title:       req.Title,
		ReleaseYear: req.ReleaseYear,
		AuthorID:    req.AuthorID,
	}
	if err := h.service.Update(book); err != nil {
		respondBookError(c, err)
		return
	}
	c.JSON(http.StatusOK, book)
}

// Delete handles DELETE /api/books/:id
func (h *BookHandler) Delete(c *gin.Context) {
	id, ok := bookIDFromParam(c)
	if !ok {
		return
	}

	if err := h.service.Delete(id); err != nil {
		respondBookError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "book deleted successfully"})
}