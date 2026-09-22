package model

// User struct represents the user model in our system
type User struct {
    ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    Name         string `gorm:"not null" json:"name"`
    Email        string `gorm:"uniqueIndex;not null" json:"email"`        //it can be duplicate for two person
    PasswordHash string `gorm:"not null" json:"password"`               
    Role         string `gorm:"default:member" json:"role"`
}