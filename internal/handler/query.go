package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/matin-sp/-online-library-management-system/internal/model"
)

// pageParams reads ?page= and ?limit=. On bad input it writes the 400
// response itself, so callers just return when ok is false.
func pageParams(c *gin.Context) (page, limit int, ok bool) {
	page, err := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(model.DefaultPage)))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page must be a positive integer"})
		return 0, 0, false
	}

	limit, err = strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(model.DefaultLimit)))
	if err != nil || limit < 1 || limit > model.MaxLimit {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a whole number between 1 and 100"})
		return 0, 0, false
	}
	return page, limit, true
}