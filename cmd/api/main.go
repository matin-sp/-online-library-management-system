package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/matin-sp/-online-library-management-system/internal/cache"
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

	// 2. Connect to Redis (optional: the app still works if Redis is down)
	redisDB := 0
	if raw := os.Getenv("REDIS_DB"); raw != "" {
		redisDB, err = strconv.Atoi(raw)
		if err != nil {
			log.Fatal("REDIS_DB must be a number: ", err)
		}
	}
	appCache := cache.New(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"), redisDB)
	defer appCache.Close()

	if err := appCache.Ping(); err != nil {
		log.Println("WARNING: Redis is not reachable, running without cache:", err)
	} else {
		fmt.Println("Successfully connected to Redis!")
	}

	// 3. Instantiate components (repository -> service -> handler)
	// Users are NOT cached: login must always read fresh data and the row
	// contains a password hash.
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// --- AUTHOR FEATURE ---
	authorRepo := repository.NewAuthorRepository(db)
	authorService := service.NewAuthorService(authorRepo, appCache)
	authorHandler := handler.NewAuthorHandler(authorService)
	// --- END AUTHOR FEATURE ---

	// --- BOOK FEATURE ---
	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo, appCache)
	bookHandler := handler.NewBookHandler(bookService)
	// --- END BOOK FEATURE ---

	// --- BOOK COPY FEATURE ---
	// Not cached: the status of a copy changes every time a book is borrowed
	// or returned, so a cached copy would show wrong data.
	bookCopyRepo := repository.NewBookCopyRepository(db)
	bookCopyService := service.NewBookCopyService(bookCopyRepo)
	bookCopyHandler := handler.NewBookCopyHandler(bookCopyService)
	// --- END BOOK COPY FEATURE ---

	// --- LOAN FEATURE ---
	// The loan service owns the database transaction, so it receives db itself.
	// The loan and fine repositories keep no db: they work on whatever
	// Queryer (db or tx) the service gives them.
	loanRepo := repository.NewLoanRepository()
	fineRepo := repository.NewFineRepository()
	loanService := service.NewLoanService(db, userRepo, bookCopyRepo, loanRepo, fineRepo)
	loanHandler := handler.NewLoanHandler(loanService)
	// --- END LOAN FEATURE ---

	// 4. Setup Gin server and define routes
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		// --- AUTHOR FEATURE: reading is public ---
		api.GET("/authors", authorHandler.GetAll)
		api.GET("/authors/:id", authorHandler.GetByID)
		// --- END AUTHOR FEATURE ---

		// --- BOOK FEATURE: reading is public ---
		api.GET("/books", bookHandler.GetAll)
		api.GET("/books/:id", bookHandler.GetByID)
		api.GET("/books/:id/with-author", bookHandler.GetByIDWithAuthor)
		// --- END BOOK FEATURE ---

		// --- BOOK COPY FEATURE: reading is public ---
		api.GET("/books/:id/copies", bookCopyHandler.ListByBook)
		// --- END BOOK COPY FEATURE ---
	}

	// Any logged-in user (valid JWT)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me/profile", userHandler.Profile)
	}

	// Admin-only routes. AuthMiddleware runs first (inherited from "protected"),
	// then RequireRole("admin") — identity first, permission second.
	admin := protected.Group("")
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.POST("/authors", authorHandler.Create)
		admin.PUT("/authors/:id", authorHandler.Update)
		admin.DELETE("/authors/:id", authorHandler.Delete)

		admin.POST("/books", bookHandler.Create)
		admin.PUT("/books/:id", bookHandler.Update)
		admin.DELETE("/books/:id", bookHandler.Delete)

		// --- BOOK COPY FEATURE: changing data needs the admin role ---
		admin.POST("/books/:id/copies", bookCopyHandler.Create)
		admin.PATCH("/copies/:id/status", bookCopyHandler.UpdateStatus)
		// --- END BOOK COPY FEATURE ---

		// --- LOAN FEATURE: the librarian checks out a copy to a user ---
		admin.POST("/loans", loanHandler.Borrow)
		// --- END LOAN FEATURE ---
	}

	router.Run(":9090")
}