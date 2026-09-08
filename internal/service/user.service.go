package service

import (
    "errors"

    "github.com/matin-sp/-online-library-management-system/internal/model"
    "github.com/matin-sp/-online-library-management-system/internal/repository"
    "golang.org/x/crypto/bcrypt"
)

// UserService handles the business logic for users
type UserService struct {
    repo *repository.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(repo *repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

// Register takes a user, hashes the password, and saves it
func (s *UserService) Register(user model.User) error {
    // 1. Hash the password using bcrypt
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
    if err != nil {
        return errors.New("failed to hash password")
    }

    // 2. Replace the plain password with the hashed password
    user.PasswordHash = string(hashedPassword)

    // 3. Ask repository to save the user
    return s.repo.Save(user)
}

// Login checks email and password, and returns user data if correct
func (s *UserService) Login(email, password string) (*model.User, error) {
    // 1. Find user by email
    user, err := s.repo.FindByEmail(email)
    if err != nil {
        return nil, errors.New("user not found")
    }

    // 2. Compare the given password with the hashed password in DB
    err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
    if err != nil {
        return nil, errors.New("wrong password")
    }

    // 3. Password is correct, return the user
    return user, nil
}