package service

import (
	"errors"
	"time"

	"github.com/matin-sp/-online-library-management-system/internal/model"
	"github.com/matin-sp/-online-library-management-system/internal/repository"
)

var (
	ErrRenewNotAllowed   = repository.ErrRenewNotAllowed
	ErrRenewLimitReached = errors.New("loan was already renewed the maximum number of times")
	ErrRenewTooLate      = errors.New("loan is past its due date and cannot be renewed")
)



// The rules (from the rules document): only the owner, only an open loan,
// only once, only before the due date, and only if the user has no unpaid fine.
// The extension is added to the existing due date.
//
// Same transaction pattern and same lock order as Borrow and Return:
// user, then loan. The user is locked so that a Return of another loan of the
// same user (which may create a fine) cannot run at the same moment.
func (s *LoanService) Renew(callerID, loanID int) (*model.Loan, error) {
	limits := borrowLimitsFor(callerID)

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	// 1. Read the loan without a lock, only to check the owner.
	loan, err := s.loans.GetByID(tx, loanID)
	if err != nil {
		return nil, err // ErrLoanNotFound or a database error
	}
	// Somebody else's loan: we answer "not found", the same as for an id that
	// does not exist. Otherwise a user could find out which loan ids exist
	// (this check is what stops the "IDOR" attack).
	if loan.UserID != callerID {
		return nil, ErrLoanNotFound
	}

	// 2. Lock the user, then the loan, and read the loan again: it may have
	// been returned while we were waiting.
	if err := s.users.LockByID(tx, callerID); err != nil {
		return nil, err
	}
	loan, err = s.loans.LockByID(tx, loanID)
	if err != nil {
		return nil, err
	}

	// 3. The rules.
	if loan.ReturnedAt != nil {
		return nil, ErrLoanAlreadyReturned
	}
	if loan.RenewalCount >= limits.maxRenewals {
		return nil, ErrRenewLimitReached
	}
	if !time.Now().Before(loan.DueDate) {
		return nil, ErrRenewTooLate
	}
	hasFine, err := s.fines.HasUnpaidByUser(tx, callerID)
	if err != nil {
		return nil, err
	}
	if hasFine {
		return nil, ErrUserHasUnpaidFines
	}
	// (The rules document also says: the copy must not be lost or damaged.
	// A copy that is checked out can never get those statuses, because
	// UpdateStatus refuses to change a checked_out copy. So there is nothing
	// to check here.)

	// 4. Move the due date and save.
	renewed, err := s.loans.Renew(tx, loanID, limits.renewalDays, limits.maxRenewals)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return renewed, nil
}