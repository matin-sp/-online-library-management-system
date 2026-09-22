package main

import (
    "fmt"
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "github.com/matin-sp/-online-library-management-system/internal/handler"
    "github.com/matin-sp/-online-library-management-system/internal/middleware"
    "github.com/matin-sp/-online-library-management-system/internal/model"
    "github.com/matin-sp/-online-library-management-system/internal/repository"
    "github.com/matin-sp/-online-library-management-system/internal/service"
)

func main() {
    // 1. Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // 2. Connect to PostgreSQL Database
    dsn := fmt.Sprintf("host=localhost user=postgres password=%s dbname=library_db port=5432 sslmode=disable", os.Getenv("DB_PASSWORD"))
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database")
    }

    // 3. Auto-Migrate (Create tables automatically)
    err = db.AutoMigrate(&model.User{})
    if err != nil {
        log.Fatal("Failed to migrate database")
    }

    // 4. Instantiate components and wire them together
    userRepo := repository.NewUserRepository(db) 
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)

    // 5. Setup Gin server and define routes
    router := gin.Default()

    api := router.Group("/api")
    {
        api.POST("/register", userHandler.Register)
        api.POST("/login", userHandler.Login)
    }

    protected := api.Group("/me")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.GET("/profile", userHandler.Profile)
    }

    // 6. Run the server on port 9090
    router.Run(":9090")
}