package model

import "time"

// Loan is one borrowing of one physical copy (a row of loans).
//
// There is no status field on purpose: the state of a loan (active, overdue,
// returned) is worked out from ReturnedAt and DueDate.
type Loan struct {
	ID           int        `json:"id"`
	UserID       int        `json:"user_id"`
	CopyID       int        `json:"copy_id"`
	BorrowedAt   time.Time  `json:"borrowed_at"`
	DueDate      time.Time  `json:"due_date"`
	ReturnedAt   *time.Time `json:"returned_at"`
	RenewalCount int        `json:"renewal_count"`
}