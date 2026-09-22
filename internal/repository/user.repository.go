package repository

import (
    "errors"

    "github.com/matin-sp/-online-library-management-system/internal/model"
    "gorm.io/gorm"
)

// UserRepository manages our database operations using GORM
type UserRepository struct {
    db *gorm.DB // Instead of a map, we now hold the database connection
}

// NewUserRepository creates a new repository with the database connection
func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

// Save adds a new user to the PostgreSQL database
func (r *UserRepository) Save(user model.User) error {
    // db.Create translates to: INSERT INTO users (...) VALUES (...)
    result := r.db.Create(&user)
    if result.Error != nil {
        // If database returns an error (like duplicate email), we send it back
        return result.Error
    }
    return nil
}

// FindByEmail looks for a user by their email in PostgreSQL
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
    var user model.User
    // db.First translates to: SELECT * FROM users WHERE email = ? LIMIT 1
    result := r.db.Where("email = ?", email).First(&user)
    if result.Error != nil {
        // gorm.ErrRecordNotFound means no user found with this email
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, errors.New("user not found")
        }
        return nil, result.Error
    }
    return &user, nil
}