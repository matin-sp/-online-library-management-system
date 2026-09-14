package main

import (
    "log"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/matin-sp/-online-library-management-system/internal/handler"
    "github.com/matin-sp/-online-library-management-system/internal/middleware"
    "github.com/matin-sp/-online-library-management-system/internal/repository"
    "github.com/matin-sp/-online-library-management-system/internal/service"
)

func main() {
    // 1. Load .env file
    // This reads the JWT_SECRET and other variables from the .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // 2. Instantiate components and wire them together (Dependency Injection)
    userRepo := repository.NewUserRepository()
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)

    // 3. Setup Gin server and define routes
    router := gin.Default()

    // Public routes (No middleware needed)
    api := router.Group("/api")
    {
        api.POST("/register", userHandler.Register)
        api.POST("/login", userHandler.Login)
    }

    // Protected routes (Middleware applied!)
    // Only users with a valid JWT token can access these
    protected := api.Group("/me")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/profile", userHandler.Profile)
    }

    // 4. Run the server on port 9090
    router.Run(":9090")
}