package middleware

import (
    "net/http"

    "github.com/gin-gonic/gin"
)


func RequireRole(allowedRoles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        value, exists := c.Get("role")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
            c.Abort()
            return
        }
 
        role, ok := value.(string)
        if !ok {
            c.JSON(http.StatusForbidden, gin.H{"error": "you do not have permission to do this"})
            c.Abort()
            return
        }

        for _, allowed := range allowedRoles {
            if role == allowed {
                c.Next()
                return
            } 
        }

        c.JSON(http.StatusForbidden, gin.H{"error": "you do not have permission to do this"})
        c.Abort()
    }
}