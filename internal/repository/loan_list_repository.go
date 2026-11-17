package repository

import (
	"fmt"

	"github.com/matin-sp/-online-library-management-system/internal/model"
)

// statusSQL is how each derived status is written in SQL. The text comes from
// THIS fixed list, never from the client, so it is safe to put it into a query.
// An empty status means "no filter".
func statusSQL(status string) (string, error) {
	switch status {
	case "":
		return "", nil
	case model.LoanStatusActive:
		return " AND l.returned_at IS NULL AND l.due_date >= now()", nil
	case model.LoanStatusOverdue:
		return " AND l.returned_at IS NULL AND l.due_date < now()", nil
	case model.LoanStatusReturned:
		return " AND l.returned_at IS NOT NULL", nil
	}
	return "", fmt.Errorf("unknown loan status %q", status)
}

// detailSelect is the SELECT part shared by both lists. The status is worked
// out here, so the database is the only place that knows what "overdue" means.
const detailSelect = `
	SELECT l.id, l.user_id, l.copy_id, l.borrowed_at, l.due_date, l.returned_at, l.renewal_count,
	       CASE WHEN l.returned_at IS NOT NULL THEN 'returned'
	            WHEN l.due_date < now() THEN 'overdue'
	            ELSE 'active' END,
	       c.book_id, b.title, c.copy_number`

const detailFrom = `
	FROM loans l
	JOIN book_copies c ON c.id = l.copy_id
	JOIN books b ON b.id = c.book_id`

// ListByUser returns one page of a user's loans, newest first.
// status is "" (all) or one of the three LoanStatus values.
func (r *LoanRepository) ListByUser(q Queryer, userID int, status string, limit, offset int) ([]model.LoanDetail, error) {
	cond, err := statusSQL(status)
	if err != nil {
		return nil, err
	}
	rows, err := q.Query(detailSelect+detailFrom+`
		WHERE l.user_id = $1`+cond+`
		ORDER BY l.borrowed_at DESC, l.id DESC
		LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loans := []model.LoanDetail{} // not nil, so the JSON shows [] and not null
	for rows.Next() {
		var d model.LoanDetail
		if err := rows.Scan(&d.ID, &d.UserID, &d.CopyID, &d.BorrowedAt, &d.DueDate, &d.ReturnedAt,
			&d.RenewalCount, &d.Status, &d.BookID, &d.BookTitle, &d.CopyNumber); err != nil {
			return nil, err
		}
		loans = append(loans, d)
	}
	return loans, rows.Err()
}

// CountByUser counts ALL the loans that match (needed for the page numbers).
func (r *LoanRepository) CountByUser(q Queryer, userID int, status string) (int, error) {
	cond, err := statusSQL(status)
	if err != nil {
		return 0, err
	}
	var n int
	err = q.QueryRow(`SELECT COUNT(*) FROM loans l WHERE l.user_id = $1`+cond, userID).Scan(&n)
	return n, err
}

// ListOverdue returns one page of the open loans that are past their due date,
// the longest overdue first. DaysOverdue is filled in by the service.
func (r *LoanRepository) ListOverdue(q Queryer, limit, offset int) ([]model.OverdueLoan, error) {
	rows, err := q.Query(detailSelect+`, u.name, u.email`+detailFrom+`
		JOIN users u ON u.id = l.user_id
		WHERE l.returned_at IS NULL AND l.due_date < now()
		ORDER BY l.due_date ASC, l.id ASC
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	loans := []model.OverdueLoan{}
	for rows.Next() {
		var o model.OverdueLoan
		d := &o.LoanDetail
		if err := rows.Scan(&d.ID, &d.UserID, &d.CopyID, &d.BorrowedAt, &d.DueDate, &d.ReturnedAt,
			&d.RenewalCount, &d.Status, &d.BookID, &d.BookTitle, &d.CopyNumber,
			&o.UserName, &o.UserEmail); err != nil {
			return nil, err
		}
		loans = append(loans, o)
	}
	return loans, rows.Err()
}

// CountOverdue counts all the overdue loans.
func (r *LoanRepository) CountOverdue(q Queryer) (int, error) {
	var n int
	err := q.QueryRow(`SELECT COUNT(*) FROM loans WHERE returned_at IS NULL AND due_date < now()`).Scan(&n)
	return n, err
}