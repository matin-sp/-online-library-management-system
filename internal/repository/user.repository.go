package repository

import (
    "errors"
    "github.com/matin-sp/-online-library-management-system/internal/model"
)

// UserRepository manages our temporary in-memory database
type UserRepository struct {
    users map[string]model.User // we use a map for faster lookup by email
}

// NewUserRepository creates a new empty repository
func NewUserRepository() *UserRepository {
    return &UserRepository{
        users: make(map[string]model.User),
    }
}

// Save adds a new user to our temporary database
func (r *UserRepository) Save(user model.User) error {
    if _, exists := r.users[user.Email]; exists {
        return errors.New("user with this email already exists")
    }
    r.users[user.Email] = user
    return nil
}

// FindByEmail looks for a user by their email
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
    user, exists := r.users[email]
    if !exists {
        return nil, errors.New("user not found")
    }
    return &user, nil
}