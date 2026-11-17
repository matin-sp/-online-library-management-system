package repository

import "github.com/matin-sp/-online-library-management-system/internal/model"

// A fine has no user_id column: it belongs to a loan, and the loan belongs to
// a user. So "the fines of a user" means "the fines of the loans of a user".

// ListByUser returns one page of a user's fines, newest first.
func (r *FineRepository) ListByUser(q Queryer, userID, limit, offset int) ([]model.Fine, error) {
	rows, err := q.Query(`
		SELECT `+fineColumns+` FROM fines
		WHERE loan_id IN (SELECT id FROM loans WHERE user_id = $1)
		ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fines := []model.Fine{}
	for rows.Next() {
		var f model.Fine
		if err := rows.Scan(&f.ID, &f.LoanID, &f.DaysLate, &f.Amount, &f.Status, &f.CreatedAt, &f.PaidAt); err != nil {
			return nil, err
		}
		fines = append(fines, f)
	}
	return fines, rows.Err()
}

// CountByUser counts all the fines of a user.
func (r *FineRepository) CountByUser(q Queryer, userID int) (int, error) {
	var n int
	err := q.QueryRow(`
		SELECT COUNT(*) FROM fines
		WHERE loan_id IN (SELECT id FROM loans WHERE user_id = $1)`, userID).Scan(&n)
	return n, err
}

// SumUnpaidByUser adds up the amount (in rial) of all unpaid fines of a user.
// COALESCE gives 0 instead of NULL when there is no unpaid fine.
func (r *FineRepository) SumUnpaidByUser(q Queryer, userID int) (int64, error) {
	var total int64
	err := q.QueryRow(`
		SELECT COALESCE(SUM(amount), 0)::bigint FROM fines
		WHERE status = 'unpaid'
		  AND loan_id IN (SELECT id FROM loans WHERE user_id = $1)`, userID).Scan(&total)
	return total, err
}