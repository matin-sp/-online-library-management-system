package repository

import (
    "database/sql"
    "errors"

    "github.com/lib/pq" 
    "github.com/matin-sp/-online-library-management-system/internal/model"
)

var (
    ErrUserNotFound    = errors.New("user not found")
    ErrDuplicateEmail  = errors.New("email already exists")
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Save(user model.User) error {
    err := r.db.QueryRow(
        "INSERT INTO users (name, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id",
        user.Name, user.Email, user.PasswordHash, user.Role,
    ).Scan(&user.ID)

    if err != nil {
        var pqErr *pq.Error
        if errors.As(err, &pqErr) && pqErr.Code == "23505" {
            return ErrDuplicateEmail
        }
        return err
    }
    return nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
    var user model.User

    err := r.db.QueryRow(
        "SELECT id, name, email, password_hash, role FROM users WHERE email = $1", email,
    ).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role)

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, ErrUserNotFound
        }
        return nil, err
    }
    return &user, nil
}