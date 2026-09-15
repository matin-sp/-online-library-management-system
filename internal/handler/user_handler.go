package handler

import (
    "errors"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "github.com/matin-sp/-online-library-management-system/internal/model"
    "github.com/matin-sp/-online-library-management-system/internal/service"
)

// UserHandler handles HTTP requests related to users
type UserHandler struct {
    service *service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(service *service.UserService) *UserHandler {
    return &UserHandler{service: service}
}

// --- DTOs (Data Transfer Objects) for Validation ---

type registerRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8,max=16"`
}

type loginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// --- Handlers ---

// Register handles POST /register
func (h *UserHandler) Register(c *gin.Context) {
    var req registerRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        var ve validator.ValidationErrors
        if errors.As(err, &ve) {
            for _, fe := range ve {
                if fe.Field() == "Password" && fe.Tag() == "max" {
                    c.JSON(http.StatusBadRequest, gin.H{"error": "entered password exceeds the allowed length"})
                    return
                }
                if fe.Field() == "Password" && fe.Tag() == "min" {
                    c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 8 characters"})
                    return
                }
                if fe.Tag() == "required" {
                    c.JSON(http.StatusBadRequest, gin.H{"error": fe.Field() + " is required"})
                    return
                }
                if fe.Tag() == "email" {
                    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
                    return
                }
            }
        }
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    user := model.User{
        Name:         req.Name,
        Email:        req.Email,
        PasswordHash: req.Password,
        Role:         "member",
    }

    if err := h.service.Register(user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

// Login handles POST /login
func (h *UserHandler) Login(c *gin.Context) {
    var req loginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        var ve validator.ValidationErrors
        if errors.As(err, &ve) {
            for _, fe := range ve {
                if fe.Tag() == "required" {
                    c.JSON(http.StatusBadRequest, gin.H{"error": fe.Field() + " is required"})
                    return
                }
                if fe.Tag() == "email" {
                    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
                    return
                }
            }
        }
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    token, err := h.service.Login(req.Email, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "login successful", "token": token})
}

// Profile handles GET /me/profile (Protected route)
func (h *UserHandler) Profile(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "user id not found in context"})
        return
    }
    role, exists := c.Get("role")
    if !exists {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "role not found in context"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Welcome to your profile!",
        "user_id": userID,
        "role":    role,
    })
}