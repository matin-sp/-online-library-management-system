package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware checks if the JWT token is valid
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Get the Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
            c.Abort()
            return
        }

        // 2. Check if the header has the "Bearer " prefix
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
            c.Abort()
            return
        }

        tokenString := parts[1]

        // 3. Parse and Validate the token
        // Note: This secret key MUST match the one in user_service.go
        secretKey := "my-super-secret-key"

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, gin.Error{Err: nil} // basic error handling
            }
            return []byte(secretKey), nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }

        // 4. Extract claims (user_id and role) and set them in the context
        // so the handler can use them later
        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            c.Set("user_id", claims["user_id"])
            c.Set("role", claims["role"])
        }

        // 5. Pass to the next handler
        c.Next()
    }
}





















