package model


import "time"

// Loan is one borrowing of one copy by one user.
// ReturnedAt is nil while the loan is active (not returned yet).
type Loan struct {
    ID           int
    UserID       int
    CopyID       int
    BorrowedAt   time.Time
    DueDate      time.Time
    ReturnedAt   *time.Time
    RenewalCount int
}