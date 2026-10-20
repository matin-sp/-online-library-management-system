package repository

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/matin-sp/-online-library-management-system/internal/model"
)

var (
	ErrAuthorNotFound = errors.New("author not found")
	ErrAuthorHasBooks = errors.New("author still has books")
	ErrAuthorExists   = errors.New("author already exists")
)

type AuthorRepository struct {
	db *sql.DB
}

func NewAuthorRepository(db *sql.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

// 23505 = unique_violation: the database refused a duplicate value.
func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}

func (r *AuthorRepository) Create(author model.Author) (*model.Author, error) {
	err := r.db.QueryRow(
		"INSERT INTO authors (name) VALUES ($1) RETURNING id",
		author.Name,
	).Scan(&author.ID)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrAuthorExists
		}
		return nil, err
	}
	return &author, nil
}

func (r *AuthorRepository) GetByID(id int) (*model.Author, error) {
	var author model.Author
	err := r.db.QueryRow(
		"SELECT id, name FROM authors WHERE id = $1", id,
	).Scan(&author.ID, &author.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAuthorNotFound
		}
		return nil, err
	}
	return &author, nil
}

func (r *AuthorRepository) GetAll() ([]model.Author, error) {
	rows, err := r.db.Query("SELECT id, name FROM authors ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	authors := make([]model.Author, 0)
	for rows.Next() {
		var a model.Author
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			return nil, err
		}
		authors = append(authors, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return authors, nil
}

func (r *AuthorRepository) Update(author model.Author) error {
	result, err := r.db.Exec(
		"UPDATE authors SET name = $1 WHERE id = $2",
		author.Name, author.ID,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrAuthorExists
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrAuthorNotFound
	}
	return nil
}

func (r *AuthorRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM authors WHERE id = $1", id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return ErrAuthorHasBooks
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrAuthorNotFound
	}
	return nil
}