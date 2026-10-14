package repository

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/matin-sp/-online-library-management-system/internal/model"
)

// Sentinel errors: named, comparable error values. Callers check them with
// errors.Is(err, repository.ErrAuthorNotFound) instead of comparing text,
// which would silently break if someone reworded the message.
var (
	ErrAuthorNotFound = errors.New("author not found")
	ErrAuthorHasBooks = errors.New("author still has books")
)

type AuthorRepository struct {
	db *sql.DB
}

func NewAuthorRepository(db *sql.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

// Create inserts a new author and returns it with the id Postgres generated.
func (r *AuthorRepository) Create(author model.Author) (*model.Author, error) {
	err := r.db.QueryRow(
		"INSERT INTO authors (name, country) VALUES ($1, $2) RETURNING id",
		author.Name, author.Country,
	).Scan(&author.ID)
	if err != nil {
		return nil, err
	}
	return &author, nil
}

// GetByID fetches a single author.
func (r *AuthorRepository) GetByID(id int) (*model.Author, error) {
	var author model.Author
	err := r.db.QueryRow(
		"SELECT id, name, country FROM authors WHERE id = $1", id,
	).Scan(&author.ID, &author.Name, &author.Country)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAuthorNotFound
		}
		return nil, err
	}
	return &author, nil
}

// GetAll fetches every author ordered by name.
func (r *AuthorRepository) GetAll() ([]model.Author, error) {
	rows, err := r.db.Query("SELECT id, name, country FROM authors ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// make(..., 0) so an empty table serializes as [] in JSON, not null.
	authors := make([]model.Author, 0)
	for rows.Next() {
		var a model.Author
		if err := rows.Scan(&a.ID, &a.Name, &a.Country); err != nil {
			return nil, err
		}
		authors = append(authors, a)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return authors, nil
}

// Update overwrites an existing author's fields.
func (r *AuthorRepository) Update(author model.Author) error {
	result, err := r.db.Exec(
		"UPDATE authors SET name = $1, country = $2 WHERE id = $3",
		author.Name, author.Country, author.ID,
	)
	if err != nil {
		return err
	}

	// An UPDATE that matches no row is not an error in SQL - it just changes
	// nothing. RowsAffected lets us report that as "not found".
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrAuthorNotFound
	}
	return nil
}

// Delete removes an author. If books still point at this author, Postgres
// refuses with a foreign key violation (SQLSTATE 23503); we translate that
// into ErrAuthorHasBooks so the API can answer 409 instead of a vague 500.
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