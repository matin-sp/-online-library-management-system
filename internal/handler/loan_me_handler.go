package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// currentUserID reads the id of the logged-in user, which AuthMiddleware put
// into the request context. A number inside a JWT is read back as float64
// (that is how JSON numbers work in Go), so we convert it to int.
// It writes the 401 response itself, so callers just return when ok is false.
func currentUserID(c *gin.Context) (int, bool) {
	raw, exists := c.Get("user_id")
	id, isNumber := raw.(float64)
	if !exists || !isNumber || id < 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return 0, false
	}
	return int(id), true
}

// Renew handles POST /api/loans/:id/renew (any logged-in user, no body).
// The user id comes from the token, so a user can only renew his own loans.
func (h *LoanHandler) Renew(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	loanID, ok := loanIDFromParam(c)
	if !ok {
		return
	}

	loan, err := h.service.Renew(userID, loanID)
	if err != nil {
		respondLoanError(c, err)
		return
	}
	c.JSON(http.StatusOK, loan)
}

// MyLoans handles GET /api/me/loans?status=&page=&limit= (any logged-in user).
func (h *LoanHandler) MyLoans(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	page, limit, ok := pageParams(c)
	if !ok {
		return
	}

	result, err := h.service.MyLoans(userID, c.Query("status"), page, limit)
	if err != nil {
		respondLoanError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// MyFines handles GET /api/me/fines?page=&limit= (any logged-in user).
func (h *LoanHandler) MyFines(c *gin.Context) {
	userID, ok := currentUserID(c)
	if !ok {
		return
	}
	page, limit, ok := pageParams(c)
	if !ok {
		return
	}

	result, err := h.service.MyFines(userID, page, limit)
	if err != nil {
		respondLoanError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"fines": result.Fines, "total_unpaid": result.TotalUnpaid})
}

// Overdue handles GET /api/loans/overdue?page=&limit= (admin only).
func (h *LoanHandler) Overdue(c *gin.Context) {
	page, limit, ok := pageParams(c)
	if !ok {
		return
	}

	result, err := h.service.OverdueLoans(page, limit)
	if err != nil {
		respondLoanError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}