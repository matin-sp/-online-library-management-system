package handler

import (
    "errors"
    "log"
    "math"
    "net/http"
    "strconv"
    "unicode/utf8"

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
// Unknown errors are logged but not shown to the client (they may contain
// internal details such as SQL text).
func respondBookError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, service.ErrBookNotFound):
        c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
    case errors.Is(err, service.ErrInvalidAuthor):
        c.JSON(http.StatusBadRequest, gin.H{"error": "author_id does not refer to an existing author"})
    case errors.Is(err, service.ErrInvalidYearRange):
        // This rule connects two fields (year_from > year_to) — only the
        // service can check it, the handler checks each field alone.
        // err.Error() returns the message "year_from must not be greater than year_to".
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

// GetAll handles GET /api/books?title=war&author_id=3&year_from=1800&year_to=1900&page=2&limit=10
//
// All query parameters are optional:
//   - title: filter books whose title contains this text (case-insensitive partial
//     match). At most 255 runes (NOT bytes — international-friendly).
//   - author_id: only books of this author (exact match). Must be a positive int.
//   - year_from: release_year >= this. Must be between 1000 and 2100 (matches
//     the database CHECK constraint books_year_valid).
//   - year_to: release_year <= this. Same range as year_from.
//   - page, limit: see pageParams.
//
// Bad values get a 400. year_from > year_to is rejected by the service
// (the rule connects two fields — only the service can check it).
// A page beyond the last one returns 200 with an empty data array.
func (h *BookHandler) GetAll(c *gin.Context) {
    page, limit, ok := pageParams(c)
    if !ok {
        return
    }

    title := c.Query("title")
    if utf8.RuneCountInString(title) > 255 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "title must be at most 255 characters"})
        return
    }

    authorID, ok := optionalIntParam(c, "author_id", 1, math.MaxInt32)
    if !ok {
        return
    }
    yearFrom, ok := optionalIntParam(c, "year_from", 1000, 2100)
    if !ok {
        return
    }
    yearTo, ok := optionalIntParam(c, "year_to", 1000, 2100)
    if !ok {
        return
    }

    result, err := h.service.List(model.BookFilter{
        Title:    title,
        AuthorID: authorID,
        YearFrom: yearFrom,
        YearTo:   yearTo,
        Page:     page,
        Limit:    limit,
    })
    if err != nil {
        respondBookError(c, err)
        return
    }
    c.JSON(http.StatusOK, result)
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