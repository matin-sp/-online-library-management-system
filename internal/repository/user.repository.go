package repository

import (
    "database/sql"
    "errors"

    "github.com/matin-sp/-online-library-management-system/internal/model"
)

type UserRepository struct {
    db *sql.DB 
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

// Save (Create) - How queries are written & safe parameters
func (r *UserRepository) Save(user model.User) error {
    // We use $1, $2 for safe parameter passing (prevents SQL Injection)
    _, err := r.db.Exec(
        "INSERT INTO users (name, email, password_hash, role) VALUES ($1, $2, $3, $4)",
        user.Name, user.Email, user.PasswordHash, user.Role,
    )
    if err != nil {
        // How database errors are handled
        return err
    }
    return nil
}

// FindByEmail (Read) - How rows are retrieved and scanned
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
    var user model.User
    
    // We write raw SQL and tell the database where to put the retrieved data using Scan()
    err := r.db.QueryRow(
        "SELECT id, name, email, password_hash, role FROM users WHERE email = $1", email,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role)
    
    // How database errors are handled (specifically "not found" error)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.New("user not found")
        }
        return nil, err
    }
    return &user, nil
}