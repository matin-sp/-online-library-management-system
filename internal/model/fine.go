package model

import "time"

// Allowed values of fines.status. They must stay the same as the CHECK
// constraint in migration 000006.
const (
	FineStatusUnpaid = "unpaid"
	FineStatusPaid   = "paid"
	FineStatusWaived = "waived"
)

// Fine is the penalty for one late return (a row of fines).
// Amount is a whole number in rial, never a decimal.
type Fine struct {
	ID        int        `json:"id"`
	LoanID    int        `json:"loan_id"`
	DaysLate  int        `json:"days_late"`
	Amount    int64      `json:"amount"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	PaidAt    *time.Time `json:"paid_at"` 
}