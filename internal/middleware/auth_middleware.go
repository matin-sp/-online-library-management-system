package middleware

import (
    "errors"
    "net/http"
    "os"
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
        // Read the secret key from environment variables (.env file)
        secretKey := os.Getenv("JWT_SECRET")

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                // Fix: Return a proper error instead of gin.Error to prevent Panic
                return nil, errors.New("unexpected signing method")
            }
            return []byte(secretKey), nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
            c.Abort()
            return
        }

        // 4. Extract claims (user_id and role) and set them in the context
        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            c.Set("user_id", claims["user_id"])
            c.Set("role", claims["role"])
        }

        // 5. Pass to the handler
        c.Next()
    }
}



















