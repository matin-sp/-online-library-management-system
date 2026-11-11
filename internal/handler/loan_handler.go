package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/matin-sp/-online-library-management-system/internal/service"
)

// LoanHandler handles HTTP requests related to loans.
type LoanHandler struct {
	service *service.LoanService
}

func NewLoanHandler(service *service.LoanService) *LoanHandler {
	return &LoanHandler{service: service}
}

// borrowRequest is the DTO for POST /api/loans.
type borrowRequest struct {
	UserID int `json:"user_id" binding:"required,gte=1"`
	CopyID int `json:"copy_id" binding:"required,gte=1"`
}

const invalidBorrowBody = "invalid request body: user_id and copy_id are required positive integers"

// respondLoanError maps service errors to HTTP status codes.
// 404 = the id does not exist. 409 = it exists, but the rules say no right now.
// Unknown errors are logged but not shown to the client.
func respondLoanError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	case errors.Is(err, service.ErrCopyNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "book copy not found"})
	case errors.Is(err, service.ErrLoanNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "loan not found"})
	case errors.Is(err, service.ErrLoanAlreadyReturned):
		c.JSON(http.StatusConflict, gin.H{"error": "this loan was already returned"})
	case errors.Is(err, service.ErrCopyNotAvailable):
		c.JSON(http.StatusConflict, gin.H{"error": "this copy is not available for borrowing"})
	case errors.Is(err, service.ErrUserHasUnpaidFines):
		c.JSON(http.StatusConflict, gin.H{"error": "user has unpaid fines"})
	case errors.Is(err, service.ErrUserHasOverdueLoans):
		c.JSON(http.StatusConflict, gin.H{"error": "user has overdue books that must be returned first"})
	case errors.Is(err, service.ErrLoanLimitReached):
		c.JSON(http.StatusConflict, gin.H{"error": "user has reached the maximum number of active loans"})
	case errors.Is(err, service.ErrAlreadyBorrowedTitle):
		c.JSON(http.StatusConflict, gin.H{"error": "user already has a copy of this book"})
	default:
		log.Println("loan handler error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

// Borrow handles POST /api/loans (admin only): the librarian checks out a copy to a user.
func (h *LoanHandler) Borrow(c *gin.Context) {
	var req borrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidBorrowBody})
		return
	}

	loan, err := h.service.Borrow(req.UserID, req.CopyID)
	if err != nil {
		respondLoanError(c, err)
		return
	}
	c.JSON(http.StatusCreated, loan)
}

// loanIDFromParam reads and validates the :id URL parameter of a loan.
// It writes the 400 response itself, so callers just return when ok is false.
func loanIDFromParam(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return 0, false
	}
	return id, true
}

// Return handles POST /api/loans/:id/return (admin only, no request body):
// the librarian checks a book back in. The answer has the closed loan and the
// fine, which is null when the book was not late.
func (h *LoanHandler) Return(c *gin.Context) {
	loanID, ok := loanIDFromParam(c)
	if !ok {
		return
	}

	result, err := h.service.Return(loanID)
	if err != nil {
		respondLoanError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"loan": result.Loan, "fine": result.Fine})
}