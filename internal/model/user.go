package model

// User struct represents the user model in our system
type User struct {
    ID           int    `json:"id"`
    Name         string `json:"name"`
    Email        string `json:"email"`
    PasswordHash string `json:"password"`
    Role         string `json:"role"` // can be "member" or "admin"
}