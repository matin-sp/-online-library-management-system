package repository

import (
	"database/sql"
	"errors"

	"github.com/matin-sp/-online-library-management-system/internal/model"
)

var (
	ErrCopyNotFound       = errors.New("book copy not found")
	ErrCopyOnLoan         = errors.New("copy is currently on loan")
	ErrCopyNumberConflict = errors.New("could not assign a copy number, please try again")
)

// maxCopyNumberAttempts is how many times Create tries again when two
// requests pick the same copy_number at the same moment.
const maxCopyNumberAttempts = 3

// copyColumns is the column list used by every query that returns a copy.
// The order must match the order of the Scan calls below.
const copyColumns = "id, book_id, copy_number, status, created_at"

type BookCopyRepository struct {
	db *sql.DB
}

func NewBookCopyRepository(db *sql.DB) *BookCopyRepository {
	return &BookCopyRepository{db: db}
}

// Create adds one new copy of a book. The copy_number is "highest number of
// this book + 1" (or 1 for the first copy), calculated inside the INSERT.
func (r *BookCopyRepository) Create(bookID int) (*model.BookCopy, error) {
	const query = `
		INSERT INTO book_copies (book_id, copy_number)
		VALUES ($1, COALESCE((SELECT MAX(copy_number) FROM book_copies WHERE book_id = $1), 0) + 1)
		RETURNING ` + copyColumns

	for attempt := 0; attempt < maxCopyNumberAttempts; attempt++ {
		var c model.BookCopy
		err := r.db.QueryRow(query, bookID).Scan(
			&c.ID, &c.BookID, &c.CopyNumber, &c.Status, &c.CreatedAt,
		)
		if err == nil {
			return &c, nil
		}
		if isForeignKeyViolation(err) {
			return nil, ErrBookNotFound // book_id does not exist in books
		}
		if isUniqueViolation(err) {
			continue // another request took this number first: calculate it again
		}
		return nil, err
	}
	return nil, ErrCopyNumberConflict
}

func (r *BookCopyRepository) GetByID(id int) (*model.BookCopy, error) {
	var c model.BookCopy
	err := r.db.QueryRow(
		"SELECT "+copyColumns+" FROM book_copies WHERE id = $1", id,
	).Scan(&c.ID, &c.BookID, &c.CopyNumber, &c.Status, &c.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCopyNotFound
		}
		return nil, err
	}
	return &c, nil
}

// ListByBook returns all copies of one book, ordered by copy_number.
// A book without copies gives an empty list; a book that does not exist
// gives ErrBookNotFound.
func (r *BookCopyRepository) ListByBook(bookID int) ([]model.BookCopy, error) {
	var bookExists bool
	err := r.db.QueryRow("SELECT EXISTS (SELECT 1 FROM books WHERE id = $1)", bookID).Scan(&bookExists)
	if err != nil {
		return nil, err
	}
	if !bookExists {
		return nil, ErrBookNotFound
	}

	rows, err := r.db.Query(
		"SELECT "+copyColumns+" FROM book_copies WHERE book_id = $1 ORDER BY copy_number", bookID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	copies := make([]model.BookCopy, 0)
	for rows.Next() {
		var c model.BookCopy
		if err := rows.Scan(&c.ID, &c.BookID, &c.CopyNumber, &c.Status, &c.CreatedAt); err != nil {
			return nil, err
		}
		copies = append(copies, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return copies, nil
}

// UpdateStatus changes the status of a copy that is NOT on loan.
// The "not on loan" check is inside the UPDATE itself (AND status <> ...),
// so there is no gap between checking and changing.
func (r *BookCopyRepository) UpdateStatus(id int, status string) (*model.BookCopy, error) {
	var c model.BookCopy
	err := r.db.QueryRow(
		"UPDATE book_copies SET status = $1 WHERE id = $2 AND status <> $3 RETURNING "+copyColumns,
		status, id, model.CopyStatusCheckedOut,
	).Scan(&c.ID, &c.BookID, &c.CopyNumber, &c.Status, &c.CreatedAt)
	if err == nil {
		return &c, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// No row was changed. Either the copy does not exist, or it is checked out.
	if _, getErr := r.GetByID(id); getErr != nil {
		return nil, getErr // ErrCopyNotFound, or a real database error
	}
	return nil, ErrCopyOnLoan
}