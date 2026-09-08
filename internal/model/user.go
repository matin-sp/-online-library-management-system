package model


// User struct represents a user in our library system

type User struct {
    ID           int    `json:"id"`
    Name         string `json:"name"`
    Email        string `json:"email"`
    PasswordHash string `json:"password"`
    Role         string `json:"role"`
}