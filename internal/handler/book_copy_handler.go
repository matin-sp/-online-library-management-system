package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/matin-sp/-online-library-management-system/internal/service"
)

// BookCopyHandler handles HTTP requests related to book copies.
type BookCopyHandler struct {
	service *service.BookCopyService
}

func NewBookCopyHandler(service *service.BookCopyService) *BookCopyHandler {
	return &BookCopyHandler{service: service}
}

// copyStatusRequest is the DTO for PATCH /api/copies/:id/status.
// "oneof" accepts only these exact words ("checked_out" is not allowed).
type copyStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=available under_repair lost damaged"`
}

const invalidCopyStatusBody = "invalid request body: status is required and must be one of: available, under_repair, lost, damaged"

// copyIDFromParam reads and validates the :id URL parameter of a copy.
// It writes the 400 response itself, so callers just return when ok is false.
func copyIDFromParam(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return 0, false
	}
	return id, true
}

// respondCopyError maps service errors to HTTP status codes.
// Unknown errors are logged but not shown to the client.
func respondCopyError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrBookNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
	case errors.Is(err, service.ErrCopyNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "book copy not found"})
	case errors.Is(err, service.ErrCopyOnLoan):
		c.JSON(http.StatusConflict, gin.H{"error": "this copy is currently on loan; its status cannot be changed by hand"})
	case errors.Is(err, service.ErrCopyNumberConflict):
		c.JSON(http.StatusConflict, gin.H{"error": "could not assign a copy number, please try again"})
	case errors.Is(err, service.ErrInvalidCopyStatus):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		log.Println("book copy handler error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func (h *BookCopyHandler) Create(c *gin.Context) {
	bookID, ok := bookIDFromParam(c)
	if !ok {
		return
	}

	bookCopy, err := h.service.Create(bookID)
	if err != nil {
		respondCopyError(c, err)
		return
	}
	c.JSON(http.StatusCreated, bookCopy)
}

func (h *BookCopyHandler) ListByBook(c *gin.Context) {
	bookID, ok := bookIDFromParam(c)
	if !ok {
		return
	}

	copies, err := h.service.ListByBook(bookID)
	if err != nil {
		respondCopyError(c, err)
		return
	}
	c.JSON(http.StatusOK, copies)
}

func (h *BookCopyHandler) UpdateStatus(c *gin.Context) {
	id, ok := copyIDFromParam(c)
	if !ok {
		return
	}

	var req copyStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidCopyStatusBody})
		return
	}

	bookCopy, err := h.service.UpdateStatus(id, req.Status)
	if err != nil {
		respondCopyError(c, err)
		return
	}
	c.JSON(http.StatusOK, bookCopy)
}
