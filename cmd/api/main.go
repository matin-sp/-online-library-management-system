package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    _ "github.com/lib/pq" // PostgreSQL driver
    "github.com/matin-sp/-online-library-management-system/internal/handler"
    "github.com/matin-sp/-online-library-management-system/internal/middleware"
    "github.com/matin-sp/-online-library-management-system/internal/repository"
    "github.com/matin-sp/-online-library-management-system/internal/service"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        log.Fatal("JWT_SECRET is missing in .env file")
    }

    // 1. Connect to PostgreSQL using standard database/sql package
    psqlInfo := fmt.Sprintf("host=localhost port=5432 user=postgres password=%s dbname=library_db sslmode=disable", os.Getenv("DB_PASSWORD"))
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal("Failed to connect to database: ", err)
    }
    defer db.Close()

    // Check if connection is actually alive
    err = db.Ping()
    if err != nil {
        log.Fatal("Failed to ping database: ", err)
    }
    fmt.Println("Successfully connected to PostgreSQL!")

    // 2. Instantiate components
    // (We will update UserRepo and add BookRepo, AuthorRepo later)
    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)

    // 3. Setup Gin server and define routes
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

    router.Run(":9090")
}