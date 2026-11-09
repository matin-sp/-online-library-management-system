package model

import "time"

const (
	CopyStatusAvailable   = "available"
	CopyStatusCheckedOut  = "checked_out"
	CopyStatusUnderRepair = "under_repair"
	CopyStatusLost        = "lost"
	CopyStatusDamaged     = "damaged"
)

type BookCopy struct {
	ID         int       `json:"id"`
	BookID     int       `json:"book_id"`
	CopyNumber int       `json:"copy_number"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
} 