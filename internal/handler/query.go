package handler

import (
    "fmt"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/matin-sp/-online-library-management-system/internal/model"
)

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


func optionalIntParam(c *gin.Context, name string, minValue, maxValue int) (value int, ok bool) {
    raw, present := c.GetQuery(name)
    if !present {
        return 0, true
    }
    value, err := strconv.Atoi(raw)
    if err != nil || value < minValue || value > maxValue {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": fmt.Sprintf("%s must be a whole number between %d and %d", name, minValue, maxValue),
        })
        return 0, false
    }
    return value, true
}