package repository

import (
	"database/sql"
	"errors"

	"github.com/matin-sp/-online-library-management-system/internal/model"
)

var (
	// ErrCopyNotAvailable means the copy cannot be lent right now. The service
	// returns it when the status is not "available", and the repository returns
	// it when the database itself refuses the loan (the partial unique index
	// uq_active_loan: the copy already has an open loan). Both mean the same to
	// the client, so both are the same error.
	ErrCopyNotAvailable = errors.New("copy is not available")

	ErrLoanNotFound        = errors.New("loan not found")
	ErrLoanAlreadyReturned = errors.New("loan is already returned")
)

// loanColumns must be in the same order as the Scan calls below.
const loanColumns = "id, user_id, copy_id, borrowed_at, due_date, returned_at, renewal_count"

// LoanRepository holds the SQL for the loans table.
// It keeps no database handle: every method receives a Queryer.
type LoanRepository struct{}

func NewLoanRepository() *LoanRepository {
	return &LoanRepository{}
}

// scanLoan reads one loan from a row. No row means the loan does not exist.
func scanLoan(row *sql.Row) (*model.Loan, error) {
	var l model.Loan
	err := row.Scan(&l.ID, &l.UserID, &l.CopyID, &l.BorrowedAt, &l.DueDate, &l.ReturnedAt, &l.RenewalCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLoanNotFound
		}
		return nil, err
	}
	return &l, nil
}

// HasOverdueByUser tells if the user still holds a book that is past its due date.
func (r *LoanRepository) HasOverdueByUser(q Queryer, userID int) (bool, error) {
	return queryExists(q, `
		SELECT EXISTS (
			SELECT 1 FROM loans
			WHERE user_id = $1 AND returned_at IS NULL AND due_date < now())`, userID)
}

// CountActiveByUser counts the loans the user has not returned yet.
func (r *LoanRepository) CountActiveByUser(q Queryer, userID int) (int, error) {
	var count int
	err := q.QueryRow(
		"SELECT COUNT(*) FROM loans WHERE user_id = $1 AND returned_at IS NULL", userID,
	).Scan(&count)
	return count, err
}

// HasActiveOfBook tells if the user already holds a copy of this book
// (any copy of the same title).
func (r *LoanRepository) HasActiveOfBook(q Queryer, userID, bookID int) (bool, error) {
	return queryExists(q, `
		SELECT EXISTS (
			SELECT 1 FROM loans l
			JOIN book_copies c ON c.id = l.copy_id
			WHERE l.user_id = $1 AND l.returned_at IS NULL AND c.book_id = $2)`, userID, bookID)
}

// Create inserts a loan. borrowed_at is now() and due_date is now() plus
// loanDays. Inside one transaction now() is a fixed time, so the two dates
// are exactly loanDays apart.
// If the partial unique index uq_active_loan refuses the row, the copy
// already has an open loan: this is the safety net behind the row lock.
func (r *LoanRepository) Create(q Queryer, userID, copyID, loanDays int) (*model.Loan, error) {
	var l model.Loan
	err := q.QueryRow(`
		INSERT INTO loans (user_id, copy_id, due_date)
		VALUES ($1, $2, now() + make_interval(days => $3))
		RETURNING `+loanColumns, userID, copyID, loanDays,
	).Scan(&l.ID, &l.UserID, &l.CopyID, &l.BorrowedAt, &l.DueDate, &l.ReturnedAt, &l.RenewalCount)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrCopyNotAvailable
		}
		return nil, err
	}
	return &l, nil
}

// GetByID reads a loan WITHOUT locking it.
func (r *LoanRepository) GetByID(q Queryer, id int) (*model.Loan, error) {
	return scanLoan(q.QueryRow("SELECT "+loanColumns+" FROM loans WHERE id = $1", id))
}

// LockByID reads a loan and locks its row until the surrounding transaction
// ends. If another transaction is returning the same loan right now, this
// call waits and then returns the NEW state (with returned_at filled in).
func (r *LoanRepository) LockByID(q Queryer, id int) (*model.Loan, error) {
	return scanLoan(q.QueryRow(
		"SELECT "+loanColumns+" FROM loans WHERE id = $1 FOR NO KEY UPDATE", id))
}

// MarkReturned closes an open loan: it sets returned_at to the real current
// time and gives back the updated loan. clock_timestamp() is the time right
// now, while now() would be the time the transaction started.
// If there is no OPEN loan with this id, it returns ErrLoanAlreadyReturned.
func (r *LoanRepository) MarkReturned(q Queryer, id int) (*model.Loan, error) {
	loan, err := scanLoan(q.QueryRow(
		"UPDATE loans SET returned_at = clock_timestamp() WHERE id = $1 AND returned_at IS NULL RETURNING "+loanColumns, id))
	if errors.Is(err, ErrLoanNotFound) {
		return nil, ErrLoanAlreadyReturned
	}
	return loan, err
}