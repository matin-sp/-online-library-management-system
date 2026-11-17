package model

import "time"


const (
	LoanStatusActive   = "active"   // with the user, not late yet
	LoanStatusOverdue  = "overdue"  // with the user, past the due date
	LoanStatusReturned = "returned" // closed
)

// IsValidLoanStatus tells if s is one of the three states above.
func IsValidLoanStatus(s string) bool {
	return s == LoanStatusActive || s == LoanStatusOverdue || s == LoanStatusReturned
}

type LoanDetail struct {
	ID           int        `json:"id"`
	UserID       int        `json:"user_id"`
	CopyID       int        `json:"copy_id"`
	BorrowedAt   time.Time  `json:"borrowed_at"`
	DueDate      time.Time  `json:"due_date"`
	ReturnedAt   *time.Time `json:"returned_at"`
	RenewalCount int        `json:"renewal_count"`
	Status       string     `json:"status"`
	BookID       int        `json:"book_id"`
	BookTitle    string     `json:"book_title"`
	CopyNumber   int        `json:"copy_number"`
}


type OverdueLoan struct {
	LoanDetail
	UserName    string `json:"user_name"`
	UserEmail   string `json:"user_email"`
	DaysOverdue int    `json:"days_overdue"`
}