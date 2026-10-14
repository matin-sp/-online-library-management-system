package repository

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/matin-sp/-online-library-management-system/internal/model"
)

var (
	ErrBookNotFound = errors.New("book not found")
	
	ErrInvalidAuthor = errors.New("author does not exist")
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}


func isForeignKeyViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23503"
}


func (r *BookRepository) Create(book model.Book) (*model.Book, error) {
	err := r.db.QueryRow(
		"INSERT INTO books (title, release_year, author_id) VALUES ($1, $2, $3) RETURNING id",
		book.Title, book.ReleaseYear, book.AuthorID,
	).Scan(&book.ID)
	if err != nil {
		if isForeignKeyViolation(err) {
			return nil, ErrInvalidAuthor
		}
		return nil, err
	}
	return &book, nil
}


func (r *BookRepository) GetByID(id int) (*model.Book, error) {
	var book model.Book
	err := r.db.QueryRow(
		"SELECT id, title, release_year, author_id FROM books WHERE id = $1", id,
	).Scan(&book.ID, &book.Title, &book.ReleaseYear, &book.AuthorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBookNotFound
		}
		return nil, err
	}
	return &book, nil
}


func (r *BookRepository) GetAll() ([]model.Book, error) {
	rows, err := r.db.Query("SELECT id, title, release_year, author_id FROM books ORDER BY title")
	if err != nil {
		return nil, err
	}
	return scanBooks(rows)
}


func (r *BookRepository) GetByAuthorID(authorID int) ([]model.Book, error) {
	rows, err := r.db.Query(
		"SELECT id, title, release_year, author_id FROM books WHERE author_id = $1 ORDER BY title",
		authorID,
	)
	if err != nil {
		return nil, err
	}
	return scanBooks(rows)
}


func scanBooks(rows *sql.Rows) ([]model.Book, error) {
	defer rows.Close()

	books := make([]model.Book, 0) 
	for rows.Next() {
		var b model.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.ReleaseYear, &b.AuthorID); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

type BookWithAuthor struct {
	Book       model.Book
	AuthorName string
}

func (r *BookRepository) GetByIDWithAuthor(id int) (*BookWithAuthor, error) {
	var result BookWithAuthor
	err := r.db.QueryRow(`
		SELECT b.id, b.title, b.release_year, b.author_id, a.name
		FROM books b
		JOIN authors a ON a.id = b.author_id
		WHERE b.id = $1`, id,
	).Scan(
		&result.Book.ID, &result.Book.Title, &result.Book.ReleaseYear, &result.Book.AuthorID,
		&result.AuthorName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBookNotFound
		}
		return nil, err
	}
	return &result, nil
}

func (r *BookRepository) Update(book model.Book) error {
	result, err := r.db.Exec(
		"UPDATE books SET title = $1, release_year = $2, author_id = $3 WHERE id = $4",
		book.Title, book.ReleaseYear, book.AuthorID, book.ID,
	)
	if err != nil {
		if isForeignKeyViolation(err) {
			return ErrInvalidAuthor
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrBookNotFound
	}
	return nil
}

func (r *BookRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrBookNotFound
	}
	return nil
}