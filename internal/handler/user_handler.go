package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
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

// Register handles POST /register
func (h *UserHandler) Register(c *gin.Context) {
    // 1. Parse and Validate the incoming JSON
    var user model.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 2. Call the service layer to do the business logic
    err := h.service.Register(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 3. Send success response
    c.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}

// Login handles POST /login
func (h *UserHandler) Login(c *gin.Context) {
    // 1. Parse the login request
    var loginData struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 2. Call the service layer to check credentials
    user, err := h.service.Login(loginData.Email, loginData.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    // 3. Send success response (we will add JWT here in the next step)
    c.JSON(http.StatusOK, gin.H{"message": "login successful", "user_id": user.ID})
}