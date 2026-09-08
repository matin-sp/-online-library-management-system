package main

import (
    "github.com/gin-gonic/gin"
    "github.com/matin-sp/-online-library-management-system/internal/handler"
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

    // Grouping public routes (no token required for now)
    api := router.Group("/api")
    {
        api.POST("/register", userHandler.Register)
        api.POST("/login", userHandler.Login)
    }

    // 3. Run the server on port 8080
    router.Run(":8080")
}