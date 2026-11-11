package service

import (
	"database/sql"
	"errors"

	"github.com/matin-sp/-online-library-management-system/internal/model"
	"github.com/matin-sp/-online-library-management-system/internal/repository"
)

// Errors the loan service can return. The first three are born in the
// repositories and only re-exported here; the others are business rules.
var (
	ErrUserNotFound     = repository.ErrUserNotFound
	ErrCopyNotAvailable = repository.ErrCopyNotAvailable

	ErrUserHasUnpaidFines   = errors.New("user has unpaid fines")
	ErrUserHasOverdueLoans  = errors.New("user has overdue loans")
	ErrLoanLimitReached     = errors.New("loan limit reached")
	ErrAlreadyBorrowedTitle = errors.New("user already has a copy of this book")
)

// LoanService owns the business rules AND the database transaction.
// Repositories only run SQL on whatever Queryer they are given.
type LoanService struct {
	db    *sql.DB
	users *repository.UserRepository
	books *repository.BookCopyRepository
	loans *repository.LoanRepository
	fines *repository.FineRepository
}

func NewLoanService(
	db *sql.DB,
	users *repository.UserRepository,
	books *repository.BookCopyRepository,
	loans *repository.LoanRepository,
	fines *repository.FineRepository,
) *LoanService {
	return &LoanService{db: db, users: users, books: books, loans: loans, fines: fines}
}

// Borrow checks out one copy to one user (the librarian does this at the desk).
// Everything happens in ONE transaction: either the loan is created AND the
// copy becomes checked_out, or nothing changes at all.
//
// RULE: after db.Begin(), every repository call gets tx and NEVER s.db.
// A call with s.db would run outside the transaction, without any error.
//
// Lock order is always: user first, then copy. This prevents deadlocks.
func (s *LoanService) Borrow(userID, copyID int) (*model.Loan, error) {
	limits := borrowLimitsFor(userID)

	// 1. Begin. The deferred Rollback undoes everything if we return early.
	// After a successful Commit it does nothing (it answers sql.ErrTxDone,
	// which we ignore on purpose).
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	// 2. Lock the user row. A second checkout of the same user waits here
	// until this transaction ends, so the count in step 5 is always fresh.
	if err := s.users.LockByID(tx, userID); err != nil {
		return nil, err // ErrUserNotFound or a database error
	}

	// 3. Unpaid fines block borrowing.
	hasFine, err := s.fines.HasUnpaidByUser(tx, userID)
	if err != nil {
		return nil, err
	}
	if hasFine {
		return nil, ErrUserHasUnpaidFines
	}

	// 4. A book that is past its due date blocks borrowing.
	hasOverdue, err := s.loans.HasOverdueByUser(tx, userID)
	if err != nil {
		return nil, err
	}
	if hasOverdue {
		return nil, ErrUserHasOverdueLoans
	}

	// 5. Loan limit. This runs after the lock, as its own statement.
	active, err := s.loans.CountActiveByUser(tx, userID)
	if err != nil {
		return nil, err
	}
	if active >= limits.maxActiveLoans {
		return nil, ErrLoanLimitReached
	}

	// 6. Lock the copy row and check that it can be lent.
	bookCopy, err := s.books.LockByID(tx, copyID)
	if err != nil {
		return nil, err // ErrCopyNotFound or a database error
	}
	if bookCopy.Status != model.CopyStatusAvailable {
		return nil, ErrCopyNotAvailable
	}

	// 7. One copy per title per user.
	hasTitle, err := s.loans.HasActiveOfBook(tx, userID, bookCopy.BookID)
	if err != nil {
		return nil, err
	}
	if hasTitle {
		return nil, ErrAlreadyBorrowedTitle
	}

	// 8. Create the loan. The partial unique index is the safety net: if it
	// refuses, the repository returns ErrCopyNotAvailable.
	loan, err := s.loans.Create(tx, userID, copyID, limits.loanDays)
	if err != nil {
		return nil, err
	}

	// 9. The copy is now out.
	if err := s.books.SetStatus(tx, copyID, model.CopyStatusCheckedOut); err != nil {
		return nil, err
	}

	// 10. Save everything.
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return loan, nil
}