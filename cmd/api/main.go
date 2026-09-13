package main

import (
    "github.com/gin-gonic/gin"
    "github.com/matin-sp/-online-library-management-system/internal/handler"
    "github.com/matin-sp/-online-library-management-system/internal/middleware"
    "github.com/matin-sp/-online-library-management-system/internal/repository"
    "github.com/matin-sp/-online-library-management-system/internal/service"
)

func main() {
    // 1. Instantiate components and wire them together (Dependency Injection)
    userRepo := repository.NewUserRepository()
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)

    // 2. Setup Gin server and define routes
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

    // 3. Run the server on port 9090
    router.Run(":9090")
}