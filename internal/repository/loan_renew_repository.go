package repository

import (
	"errors"

	"github.com/matin-sp/-online-library-management-system/internal/model"
)


var ErrRenewNotAllowed = errors.New("loan cannot be renewed")

func (r *LoanRepository) Renew(q Queryer, id, extraDays, maxRenewals int) (*model.Loan, error) {
	l, err := scanLoan(q.QueryRow(`
		UPDATE loans
		SET due_date = due_date + make_interval(days => $2),
		    renewal_count = renewal_count + 1
		WHERE id = $1
		  AND returned_at IS NULL
		  AND renewal_count < $3
		  AND due_date > clock_timestamp()
		RETURNING `+loanColumns, id, extraDays, maxRenewals))
	if errors.Is(err, ErrLoanNotFound) {
		return nil, ErrRenewNotAllowed
	}
	return l, err
}