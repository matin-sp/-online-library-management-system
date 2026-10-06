package model

type User struct {
    ID           uint    `json:"id"`
    Name         string `json:"name"`
    Email        string `json:"email"`
    PasswordHash string `json:"-"`
    Role         string `json:"role"`
}