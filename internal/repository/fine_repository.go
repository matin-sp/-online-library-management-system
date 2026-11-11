package repository

import "github.com/matin-sp/-online-library-management-system/internal/model"

// fineColumns must be in the same order as the Scan call in Create.
const fineColumns = "id, loan_id, days_late, amount, status, created_at, paid_at"

// FineRepository holds the SQL for the fines table.
// It keeps no database handle: every method receives a Queryer.
type FineRepository struct{}

func NewFineRepository() *FineRepository {
	return &FineRepository{}
}

// HasUnpaidByUser tells if the user has an unpaid fine from ANY past loan.
func (r *FineRepository) HasUnpaidByUser(q Queryer, userID int) (bool, error) {
	return queryExists(q, `
		SELECT EXISTS (
			SELECT 1 FROM fines f
			JOIN loans l ON l.id = f.loan_id
			WHERE l.user_id = $1 AND f.status = 'unpaid')`, userID)
}

// Create inserts the fine of one loan (status 'unpaid' is the table default).
// loan_id is UNIQUE, so a loan can never get two fines. While the loan row is
// locked and checked this cannot happen; if it ever does, it is a bug, so the
// raw database error is returned and the handler answers 500.
func (r *FineRepository) Create(q Queryer, loanID, daysLate int, amount int64) (*model.Fine, error) {
	var f model.Fine
	err := q.QueryRow(`
		INSERT INTO fines (loan_id, days_late, amount)
		VALUES ($1, $2, $3)
		RETURNING `+fineColumns, loanID, daysLate, amount,
	).Scan(&f.ID, &f.LoanID, &f.DaysLate, &f.Amount, &f.Status, &f.CreatedAt, &f.PaidAt)
	if err != nil {
		return nil, err
	}
	return &f, nil
}