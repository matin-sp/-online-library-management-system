package handler

import (
    "errors"
    "log"
    "net/http"
    "strconv"
    "unicode/utf8"

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


type authorRequest struct {
    Name string `json:"name" binding:"required,max=100"`
}

const invalidAuthorBody = "invalid request body: name is required (max 100 characters)"

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

    author, err := h.service.Create(model.Author{Name: req.Name})
    if err != nil {
        respondAuthorError(c, err)
        return
    }
    c.JSON(http.StatusCreated, author)
}

// GetAll handles GET /api/authors and GET /api/authors?name=dost&page=2&limit=10
//
// Query parameters (all optional):
//   - name: filter authors whose name contains this text (case-insensitive
//     partial match). At most 100 runes (NOT 100 bytes - international-friendly).
//   - page: which page to return. Must be >= 1. Defaults to 1.
//   - limit: how many rows per page. Must be between 1 and 100. Defaults to 10.
//
// Bad values (e.g. ?page=0 or ?limit=1000) return 400 Bad Request.
// A page beyond the last one (e.g. ?page=99 with only 2 pages of data) returns
// 200 with an empty data array - that is normal, not an error.
func (h *AuthorHandler) GetAll(c *gin.Context) {
    // 1) Parse and validate pagination parameters (?page=, ?limit=).
    //    pageParams writes the 400 response itself if input is bad, so we
    //    just return here.
    page, limit, ok := pageParams(c)
    if !ok {
        return
    }

    // 2) Read the optional ?name= filter. c.Query returns "" if absent.
    //    We validate the length in runes (not bytes) so that Persian, Arabic,
    //    Chinese, etc. names get the same 100-character limit as English ones.
    name := c.Query("name")
    if utf8.RuneCountInString(name) > 100 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "name must be at most 100 characters"})
        return
    }

    // 3) Build the filter and ask the service for one page of results.
    //    The service will apply TrimSpace(name) and clamp page/limit to safe
    //    defaults (defense in depth) before calling the repository.
    result, err := h.service.List(model.AuthorFilter{Name: name, Page: page, Limit: limit})
    if err != nil {
        respondAuthorError(c, err)
        return
    }

    // 4) result is a *model.Page[model.Author], which serializes to:
    //        { "data": [...], "page": 2, "limit": 10, "total": 14, "total_pages": 2 }
    //    Note: the response shape changed from a plain array to an object
    //    with pagination metadata. The handler does not need to do anything
    //    special - we just hand the Page struct to gin's JSON serializer.
    c.JSON(http.StatusOK, result)
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

// Update handles PUT /api/authors/:id (full replacement of name)
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

    author := model.Author{ID: id, Name: req.Name}
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