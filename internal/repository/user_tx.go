package repository

import (
	"database/sql"
	"errors"
)

// LockByID locks the user's row until the surrounding transaction ends.
// Another transaction that wants the same row waits here. The checkout flow
// uses this to run the checkouts of ONE user one after the other, so two
// requests cannot both pass the loan-limit check.
//
// FOR NO KEY UPDATE is the weaker lock: it does not block other rows that
// only point to this user (like a new loan).
func (r *UserRepository) LockByID(q Queryer, id int) error {
	var found int
	err := q.QueryRow("SELECT id FROM users WHERE id = $1 FOR NO KEY UPDATE", id).Scan(&found)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrUserNotFound
	}
	return err
}